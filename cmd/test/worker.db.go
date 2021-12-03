package main

import (
	"github.com/Myriad-Dreamin/aliali/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func (w *Worker) syncDBSchemas(db *gorm.DB) *gorm.DB {
	w.s.Suppress(db.AutoMigrate(&model.AliAuthModel{}))
	w.s.Suppress(db.AutoMigrate(&model.UploadModel{}))
	return db
}

func (w *Worker) openDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(w.dbPath))
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
