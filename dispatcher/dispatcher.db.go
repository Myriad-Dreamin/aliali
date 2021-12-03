package dispatcher

import (
	"github.com/Myriad-Dreamin/aliali/model"
	"github.com/Myriad-Dreamin/aliali/pkg/suppress"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	logger_util "gorm.io/gorm/logger"
)

type DBManager struct {
	S suppress.ISuppress
}

func (d *DBManager) syncDBSchemas(db *gorm.DB) *gorm.DB {
	d.S.Suppress(db.AutoMigrate(&model.AliAuthModel{}))
	d.S.Suppress(db.AutoMigrate(&model.UploadModel{}))
	return db
}

func (d *DBManager) OpenSqliteDB(dsn string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger_util.Default.LogMode(logger_util.Silent),
	})
	if err != nil {
		d.S.Suppress(err)
		return nil
	}

	return d.syncDBSchemas(db)
}

func (d *Dispatcher) syncDBSchemas(db *gorm.DB) *gorm.DB {
	d.s.Suppress(db.AutoMigrate(&model.AliAuthModel{}))
	d.s.Suppress(db.AutoMigrate(&model.UploadModel{}))
	return db
}

func (d *Dispatcher) openDB() *gorm.DB {
	return d.dbMgr.OpenSqliteDB(d.dbPath)
}

func (d *Dispatcher) openMock() *gorm.DB {
	return d.dbMgr.OpenSqliteDB(":memory:")
}

func (d *Dispatcher) deleteAllUploadData() {
	d.db.Debug().Model(&model.UploadModel{}).Where("id != ?", 0).Delete(&model.UploadModel{})
}
