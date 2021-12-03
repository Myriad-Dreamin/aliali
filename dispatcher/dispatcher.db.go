package dispatcher

import (
	"github.com/Myriad-Dreamin/aliali/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	logger_util "gorm.io/gorm/logger"
)

func (d *Dispatcher) syncDBSchemas(db *gorm.DB) *gorm.DB {
	d.s.Suppress(db.AutoMigrate(&model.AliAuthModel{}))
	d.s.Suppress(db.AutoMigrate(&model.UploadModel{}))
	return db
}

func (d *Dispatcher) openDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(d.dbPath), &gorm.Config{
		Logger: logger_util.Default.LogMode(logger_util.Silent),
	})
	if err != nil {
		d.s.Suppress(err)
		return nil
	}

	return d.syncDBSchemas(db)
}

func (d *Dispatcher) openMock() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"))
	if err != nil {
		d.s.Suppress(err)
		return nil
	}

	return d.syncDBSchemas(db)
}

func (d *Dispatcher) deleteAllUploadData() {
	d.db.Debug().Model(&model.UploadModel{}).Where("id != ?", 0).Delete(&model.UploadModel{})
}
