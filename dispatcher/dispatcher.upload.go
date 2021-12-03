package dispatcher

import (
	"context"
	"github.com/Myriad-Dreamin/aliali/dispatcher/service"
	"github.com/Myriad-Dreamin/aliali/model"
	ali_drive "github.com/Myriad-Dreamin/aliali/pkg/ali-drive"
	"github.com/Myriad-Dreamin/aliali/pkg/ali-notifier"
	"os"
)

func (d *Dispatcher) serveUploadRequest(
	ifs FsClearInterface, req *ali_notifier.FsUploadRequest, uploadReq service.IUploadRequest) error {
	if !d.ensureFsFileExists(ifs, req.LocalPath) {
		return nil
	}
	if !d.xdb.TransitUploadStatus(d.db, req, model.UploadStatusInitialized, model.UploadStatusUploading) {
		return nil
	}

	serviceImpl := <-d.serviceQueue

	go func() {
		resp := serviceImpl.Upload(context.TODO(), d.authedAli, uploadReq)

		var returnService = func() {
			// the service is blameless
			d.serviceQueue <- serviceImpl
		}

		// handle error after returning service
		// we leverage the bottom half of this coroutine to process an upload response
		switch resp.Code {
		case service.UploadOK:
			if !d.ensureFsFileExists(ifs, req.LocalPath) {
				returnService()
				return
			}
			if d.xdb.TransitUploadStatus(d.db, req, model.UploadStatusUploading, model.UploadStatusUploaded) {
				d.checkUploadAndClear(ifs, req)
			}
		case service.UploadCancelled:
			var m = &model.UploadModel{
				ID:         req.TransactionID,
				DriveID:    req.DriveID,
				RemotePath: req.RemotePath,
				LocalPath:  req.LocalPath,
			}

			if m.Set(d.s, uploadReq.Session()) {
				d.xdb.SaveUploadSession(d.db, m)
			}
			d.xdb.TransitUploadStatus(d.db, req, model.UploadStatusUploading, model.UploadStatusInitialized)
		default:
			if resp.Err != nil {
				d.s.Suppress(resp.Err)
			}
			d.xdb.TransitUploadStatus(d.db, req, model.UploadStatusUploading, model.UploadStatusInitialized)
		}
		returnService()
	}()
	return nil
}

func (d *Dispatcher) serveFsUploadRequest(req *ali_notifier.FsUploadRequest) error {
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
					DriveID: d.cfg.AliDrive.DriveId,
				},
				PartInfoList: nil,
				UploadID:     "",
			},
			XFileName:  req.RemotePath,
			XSize:      st.Size(),
			XChunkHint: d.chunkSize(),
		},
		R: o,
		S: d.s,
	}

	return d.serveUploadRequest(realFs(), req, uploadReq)
}