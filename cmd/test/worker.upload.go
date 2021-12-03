package main

import (
	"context"
	"github.com/Myriad-Dreamin/aliali/model"
	"github.com/Myriad-Dreamin/aliali/pkg/ali-notifier"
	"os"
)

func (w *Worker) serveUploadRequest(
	ifs FsClearInterface, req *ali_notifier.FsUploadRequest, uploadReq IUploadRequest) error {
	if !w.xdb.TransitUploadStatus(w.db, ifs, req, model.UploadStatusInitialized, model.UploadStatusUploading) {
		return nil
	}

	workerImpl := <-w.serviceQueue

	go func() {
		resp := workerImpl.Upload(context.TODO(), w.authedAli, uploadReq)

		// the worker is blameless
		w.serviceQueue <- workerImpl

		// handle error after returning worker
		// we leverage the bottom half of this coroutine to process an upload response
		if resp.Code == 0 {
			if w.xdb.TransitUploadStatus(w.db, ifs, req, model.UploadStatusUploading, model.UploadStatusUploaded) {
				w.checkUploadAndClear(ifs, req)
			}
		} else if resp.Err != nil {
			w.s.Suppress(resp.Err)
		}
	}()
	return nil
}

func (w *Worker) serveFsUploadRequest(req *ali_notifier.FsUploadRequest) error {
	o, err := os.OpenFile(req.LocalPath, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}

	st, err := o.Stat()
	if err != nil {
		return err
	}

	uploadReq := &RandReaderUploadRequest{
		BaseUploadRequest: BaseUploadRequest{
			// multiple drive?
			XDriverID:  w.cfg.AliDrive.DriveId,
			XFileName:  req.RemotePath,
			XSize:      st.Size(),
			XChunkHint: w.chunkSize(),
		},
		r: o,
		s: w.s,
	}

	return w.serveUploadRequest(realFs(), req, uploadReq)
}
