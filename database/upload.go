package database

import (
	"database/sql"
	"github.com/Myriad-Dreamin/aliali/model"
	"github.com/Myriad-Dreamin/aliali/pkg/ali-notifier"
	"gorm.io/gorm"
	"io/fs"
	"os"
)

func (s *DB) TransitUploadStatus(
	db *gorm.DB, ifs fs.FS, req *ali_notifier.FsUploadRequest, fr, to int) bool {
	return s.TransitUploadStatusT(db, ifs, req, func(req *ali_notifier.FsUploadRequest, status int) (target int, err error) {
		if status != fr {
			return
		}

		target = to
		return
	})
}

func (s *DB) TransitUploadStatusT(
	db *gorm.DB, operating fs.FS, req *ali_notifier.FsUploadRequest,
	transition func(req *ali_notifier.FsUploadRequest, status int) (int, error)) (changed bool) {

	// return anyway
	if _, err := fs.Stat(operating, req.LocalPath); os.IsNotExist(err) {
		return
		// fs error
	} else if err != nil && !os.IsExist(err) {
		s.Suppress(err)
		return
	}

	if err := db.Transaction(func(tx *gorm.DB) error {
		// querying
		var x = model.UploadModel{
			ID: req.TransactionID,
		}

		// read
		e := tx.Model(&x).Select("status").First(&x)
		if e.Error == gorm.ErrRecordNotFound {
			s.WarnOnce(e.Error)
			return nil
		} else if e.Error != nil {
			return e.Error
		}

		// transit
		var st, err = transition(req, x.Status)
		if err != nil {
			return err
		}
		if st == 0 {
			return nil
		}
		changed = true

		// commit
		if x.Status != x.Status {
			e = tx.Model(&x).Update("status", st)
			if e.Error == gorm.ErrRecordNotFound {
				s.WarnOnce(e.Error)
				return nil
			} else if e.Error != nil {
				tx.Rollback()
				changed = false
				return e.Error
			}
		}

		return nil
	}, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	}); err != nil {
		s.Suppress(err)
		return
	}

	return
}
