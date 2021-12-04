package dispatcher

import (
	"github.com/Myriad-Dreamin/aliali/model"
	ali_notifier "github.com/Myriad-Dreamin/aliali/pkg/ali-notifier"
	"time"
)

func (d *Dispatcher) Loop() {
	d.logger.Printf("main dispatcher entering forever loop")
	tick := time.NewTicker(time.Second)
	var stackModel model.UploadModel

	var getReq = func() *ali_notifier.FsUploadRequest {
		return &ali_notifier.FsUploadRequest{
			TransactionID: stackModel.ID,
			DriveID:       stackModel.DriveID,
			RemotePath:    stackModel.RemotePath,
			LocalPath:     stackModel.LocalPath,
		}
	}

	for {
		if !d.xdb.FindMatchedStatusRequest(d.db, &stackModel, model.UploadStatusUploading) {
			break
		}
		d.xdb.TransitUploadStatus(d.db, getReq(), model.UploadStatusUploading, model.UploadStatusInitialized)
	}

	for {
		if !d.xdb.FindMatchedStatusRequest(d.db, &stackModel, model.UploadStatusUploaded) {
			break
		}
		d.xdb.TransitUploadStatus(d.db, getReq(), model.UploadStatusUploaded, model.UploadStatusInitialized)
	}

	for {
		select {
		case req := <-d.fileUploads:
			d.logger.Printf("receive file upload request: %v", req)
			if d.authExpired() {
				d.refreshAuth()
			}
			if err := d.serveFsUploadRequest(req); err != nil {
				d.s.Suppress(err)
			}
		case <-tick.C:
			if d.xdb.FindMatchedStatusRequest(d.db, &stackModel, model.UploadStatusInitialized) {
				d.fileUploads <- getReq()
			}
		}
	}
}
