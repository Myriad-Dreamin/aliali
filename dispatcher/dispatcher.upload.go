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
		d.xdb.TransitUploadStatus(d.db, req, model.UploadStatusInitialized, model.UploadStatusSettledExitFileFlyAway)
		return nil
	}
	if !d.xdb.TransitUploadStatus(d.db, req, model.UploadStatusInitialized, model.UploadStatusUploading) {
		return nil
	}

	serviceCtx := d.waitService()

	go func() {
		resp := serviceCtx.Impl.Upload(context.TODO(), serviceCtx.authedAli, uploadReq)

		var returnService = func() {
			// the service is blameless
			d.serviceQueue <- serviceCtx
		}

		var saveSession = func() {
			var m = &model.UploadModel{
				ID:         req.TransactionID,
				DriveID:    req.DriveID,
				RemotePath: req.RemotePath,
				LocalPath:  req.LocalPath,
			}

			if m.Set(d.s, uploadReq.Session()) {
				d.xdb.SaveUploadSession(d.db, m)
			}
		}

		if resp == nil {
			d.xdb.TransitUploadStatus(d.db, req, model.UploadStatusUploading, model.UploadStatusInitialized)
			returnService()
			return
		}

		// handle error after returning service
		// we leverage the bottom half of this coroutine to process an upload response
		switch resp.Code {
		case service.UploadOK:
			saveSession()
			if !d.ensureFsFileExists(ifs, req.LocalPath) {
				d.xdb.TransitUploadStatus(d.db, req, model.UploadStatusUploading, model.UploadStatusSettledExitFileFlyAway)
				returnService()
				return
			}
			if d.xdb.TransitUploadStatus(d.db, req, model.UploadStatusUploading, model.UploadStatusUploaded) {
				d.checkUploadAndClear(ifs, req)
			}
			returnService()
			return
		case service.UploadCancelled:
			saveSession()
			d.xdb.TransitUploadStatus(d.db, req, model.UploadStatusUploading, model.UploadStatusInitialized)
			returnService()
			return
		default:
			if resp.Err != nil {
				d.s.Suppress(resp.Err)
			}
			d.xdb.TransitUploadStatus(d.db, req, model.UploadStatusUploading, model.UploadStatusInitialized)
			returnService()
			return
		}
	}()
	return nil
}

func (d *Dispatcher) serveFsUploadRequest(req *ali_notifier.FsUploadRequest) error {
	var operating = realFs()
	if !d.ensureFsFileExists(operating, req.LocalPath) {
		d.xdb.TransitUploadStatus(d.db, req, model.UploadStatusInitialized, model.UploadStatusSettledExitFileFlyAway)
		return nil
	}

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

	return d.serveUploadRequest(operating, req, uploadReq)
}
