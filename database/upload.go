package database

import (
	"database/sql"
	"github.com/Myriad-Dreamin/aliali/model"
	"github.com/Myriad-Dreamin/aliali/pkg/ali-notifier"
	"gorm.io/gorm"
)

func (s *DB) FindUploadRequest(db *gorm.DB, req *model.UploadModel) bool {
	e := db.Where(model.UploadModel{
		DriveID:    req.DriveID,
		RemotePath: req.RemotePath,
		LocalPath:  req.LocalPath,
	}).First(req)
	if e.Error == gorm.ErrRecordNotFound {
		return false
	} else if e.Error != nil {
		s.Suppress(e.Error)
		return false
	}

	return true
}

func (s *DB) FindMatchedStatusRequest(db *gorm.DB, req *model.UploadModel, st int) bool {
	req.ID = 0
	e := db.Model(req).Where("status = ?", st).First(req)
	if e.Error == gorm.ErrRecordNotFound {
		return false
	} else if e.Error != nil {
		s.Suppress(e.Error)
		return false
	}

	return true
}

func (s *DB) SaveUploadRequest(db *gorm.DB, model *model.UploadModel) {
	if model.ID == 0 {
		if !s.FindUploadRequest(db, model) {
			db = db.Create(model)
		} else {
			db = db.Save(model)
		}
	} else {
		db = db.Save(model)
	}
	if db.Error != nil {
		s.Suppress(db.Error)
		return
	}

	return
}

func (s *DB) SaveUploadSession(db *gorm.DB, model *model.UploadModel) bool {
	if model.ID == 0 {
		if !s.FindUploadRequest(db, model) {
			return false
		} else {
			db = db.Model(model).Update("raw", model.Raw)
		}
	} else {
		db = db.Model(model).Update("raw", model.Raw)
	}
	if db.Error != nil {
		s.Suppress(db.Error)
		return false
	}

	return true
}

func (s *DB) TransitUploadStatus(
	db *gorm.DB, req *ali_notifier.FsUploadRequest, fr, to int) bool {
	return s.TransitUploadStatusT(db, req, func(req *ali_notifier.FsUploadRequest, status int) (target int, err error) {
		if status != fr {
			return
		}

		target = to
		return
	})
}

func (s *DB) TransitUploadStatusT(
	db *gorm.DB, req *ali_notifier.FsUploadRequest,
	transition func(req *ali_notifier.FsUploadRequest, status int) (int, error)) (changed bool) {
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
		if x.Status != st {
			e = tx.Model(&x).Debug().Update("status", st)
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
