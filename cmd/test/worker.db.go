package main

import (
	"database/sql"
	"github.com/Myriad-Dreamin/aliali/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"io/fs"
	"os"
)

func (w *Worker) syncDBSchemas(db *gorm.DB) *gorm.DB {
	w.s.Suppress(db.AutoMigrate(&model.AliAuthModel{}))
	w.s.Suppress(db.AutoMigrate(&model.UploadModel{}))
	return db
}

func (w *Worker) openDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("ali.db"))
	if err != nil {
		w.s.Suppress(err)
		return nil
	}

	return w.syncDBSchemas(db)
}

func (w *Worker) openMock() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"))
	if err != nil {
		w.s.Suppress(err)
		return nil
	}

	return w.syncDBSchemas(db)
}

func (w *Worker) transitUploadStatus(ifs fs.FS, req *FsUploadRequest, fr, to int) bool {
	return w.fileUploadStatusTransition(ifs, req, func(req *FsUploadRequest, status int) (target int, err error) {
		if status != fr {
			return
		}

		target = to
		return
	})
}

func (w *Worker) fileUploadStatusTransition(
	operating fs.FS, req *FsUploadRequest, transition func(req *FsUploadRequest, status int) (int, error)) (changed bool) {

	// return anyway
	if _, err := fs.Stat(operating, req.LocalPath); os.IsNotExist(err) {
		return
		// fs error
	} else if err != nil && !os.IsExist(err) {
		w.s.Suppress(err)
		return
	}

	if err := w.db.Transaction(func(tx *gorm.DB) error {
		// querying
		var x = model.UploadModel{
			ID: req.TransactionID,
		}

		// read
		e := tx.Model(&x).Select("status").First(&x)
		if e.Error == gorm.ErrRecordNotFound {
			w.warnOnce(e.Error)
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
				w.warnOnce(e.Error)
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
		w.s.Suppress(err)
		return
	}

	return
}
