package dispatcher

import (
	"github.com/Myriad-Dreamin/aliali/model"
	ali_notifier "github.com/Myriad-Dreamin/aliali/pkg/ali-notifier"
)

func (d *Dispatcher) setupNotifier() {
	d.notifier.RegisterCallback(d)
}

func (d *Dispatcher) OnFsUpload(req *ali_notifier.FsUploadRequest) {
	var m = fsUploadReq2Model(req)

	if len(m.DriveID) == 0 {
		m.DriveID = d.cfg.AliDrive.DriveId
		req.DriveID = m.DriveID
	}

	if !d.xdb.FindUploadRequest(d.db, m) {
		// todo: check it
		m.Raw = []byte("{}")
		d.xdb.SaveUploadRequest(d.db, m)
		req.TransactionID = m.ID
		if d.xdb.TransitUploadStatus(d.db, req, model.UploadStatusUninitialized, model.UploadStatusInitialized) {
			d.fileUploads <- req
		}
		return
	}
}
