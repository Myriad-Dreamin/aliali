package main

import (
	"context"
	"github.com/Myriad-Dreamin/aliali/model"
	ali_drive "github.com/Myriad-Dreamin/aliali/pkg/ali-drive"
	"github.com/Myriad-Dreamin/aliali/pkg/ali-notifier"
	"os"
)

func (w *Worker) serveUploadRequest(
	ifs FsClearInterface, req *ali_notifier.FsUploadRequest, uploadReq IUploadRequest) error {
	if !w.ensureFsFileExists(ifs, req.LocalPath) {
		return nil
	}
	if !w.xdb.TransitUploadStatus(w.db, req, model.UploadStatusInitialized, model.UploadStatusUploading) {
		return nil
	}

	workerImpl := <-w.serviceQueue

	go func() {
		resp := workerImpl.Upload(context.TODO(), w.authedAli, uploadReq)

		var returnService = func() {
			// the worker is blameless
			w.serviceQueue <- workerImpl
		}

		// handle error after returning worker
		// we leverage the bottom half of this coroutine to process an upload response
		switch resp.Code {
		case UploadOK:
			if !w.ensureFsFileExists(ifs, req.LocalPath) {
				returnService()
				return
			}
			if w.xdb.TransitUploadStatus(w.db, req, model.UploadStatusUploading, model.UploadStatusUploaded) {
				w.checkUploadAndClear(ifs, req)
			}
		case UploadCancelled:
			var m = &model.UploadModel{
				ID:         req.TransactionID,
				DriveID:    req.DriveID,
				RemotePath: req.RemotePath,
				LocalPath:  req.LocalPath,
			}

			if m.Set(w.s, uploadReq.Session()) {
				w.xdb.SaveUploadSession(w.db, m)
			}
		default:
			if resp.Err != nil {
				w.s.Suppress(resp.Err)
			}
		}
		returnService()
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
			XSession: &ali_drive.UploadSession{
				DriveDirentID: ali_drive.DriveDirentID{
					DriveID: w.cfg.AliDrive.DriveId,
				},
				PartInfoList: nil,
				UploadID:     "",
			},
			XFileName:  req.RemotePath,
			XSize:      st.Size(),
			XChunkHint: w.chunkSize(),
		},
		r: o,
		s: w.s,
	}

	return w.serveUploadRequest(realFs(), req, uploadReq)
}
