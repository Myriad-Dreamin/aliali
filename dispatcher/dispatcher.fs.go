package dispatcher

import (
	"github.com/Myriad-Dreamin/aliali/model"
	"github.com/Myriad-Dreamin/aliali/pkg/ali-notifier"
	"io/fs"
	"os"
	"strings"
)

type FsClearInterface interface {
	fs.FS
	Remove(path string) error
}

func (d *Dispatcher) ensureFsFileExists(operating FsClearInterface, path string) bool {
	if _, err := fs.Stat(operating, strings.TrimPrefix(path, "/")); os.IsNotExist(err) {
		return false
		// fs error
	} else if err != nil && !os.IsExist(err) {
		d.s.Suppress(err)
		return false
	}

	return true
}

func (d *Dispatcher) checkUploadAndClear(operating FsClearInterface, req *ali_notifier.FsUploadRequest) {
	if !d.ensureFsFileExists(operating, req.LocalPath) {
		return
	}
	d.xdb.TransitUploadStatusT(d.db, req, func(req *ali_notifier.FsUploadRequest, status int) (targetStatus int, e error) {
		if status != model.UploadStatusUploaded {
			return
		}

		targetStatus = model.UploadStatusSettledClear
		// return anyway
		if _, err := fs.Stat(operating, strings.TrimPrefix(req.LocalPath, "/")); os.IsNotExist(err) {
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
