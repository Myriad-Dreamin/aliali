package main

import (
	"fmt"
	"github.com/Myriad-Dreamin/aliali/model"
	ali_notifier "github.com/Myriad-Dreamin/aliali/pkg/ali-notifier"
)

func (w *Worker) setupNotifier() {
	w.notifier.RegisterCallback(w)
}

func (w *Worker) OnFsUpload(req *ali_notifier.FsUploadRequest) {
	fmt.Println(req)
	var m = &model.UploadModel{
		DriveID:    req.DriveID,
		RemotePath: req.RemotePath,
		LocalPath:  req.LocalPath,
		Raw:        nil,
	}

	if !w.xdb.FindUploadRequest(w.db, m) {
		// todo: check it
		m.Raw = []byte("{}")
		w.xdb.SaveUploadRequest(w.db, m)
		w.xdb.TransitUploadStatus(w.db, req, model.UploadStatusUninitialized, model.UploadStatusInitialized)
		return
	}
}
