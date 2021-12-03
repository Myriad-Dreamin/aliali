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

func (w *Worker) ensureFsFileExists(operating FsClearInterface, path string) bool {
	if _, err := fs.Stat(operating, path); os.IsNotExist(err) {
		return false
		// fs error
	} else if err != nil && !os.IsExist(err) {
		w.s.Suppress(err)
		return false
	}

	return true
}

func (w *Worker) checkUploadAndClear(operating FsClearInterface, req *ali_notifier.FsUploadRequest) {
	if !w.ensureFsFileExists(operating, req.LocalPath) {
		return
	}
	w.xdb.TransitUploadStatusT(w.db, req, func(req *ali_notifier.FsUploadRequest, status int) (targetStatus int, e error) {
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
