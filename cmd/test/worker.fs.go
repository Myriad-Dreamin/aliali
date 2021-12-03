package main

import (
	"github.com/Myriad-Dreamin/aliali/model"
	"github.com/Myriad-Dreamin/aliali/pkg/ali-notifier"
	"io/fs"
	"os"
)

type FsClearInterface interface {
	fs.FS
	Remove(path string) error
}

func (w *Worker) checkUploadAndClear(operating FsClearInterface, req *ali_notifier.FsUploadRequest) {
	w.xdb.TransitUploadStatusT(w.db, operating, req, func(req *ali_notifier.FsUploadRequest, status int) (targetStatus int, e error) {
		if status != model.UploadStatusUploaded {
			return
		}

		targetStatus = model.UploadStatusSettledClear
		// return anyway
		if _, err := fs.Stat(operating, req.LocalPath); os.IsNotExist(err) {
			return
			// fs error
		} else if err != nil && !os.IsExist(err) {
			e = err
			return
		}

		if err := operating.Remove(req.LocalPath); err != nil {
			e = err
			return
		}
		targetStatus = model.UploadStatusSettledClear
		return
	})
}
