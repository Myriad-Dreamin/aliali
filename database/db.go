package database

import (
  "github.com/Myriad-Dreamin/aliali/models"
  "github.com/Myriad-Dreamin/aliali/pkg/suppress"
  "gorm.io/gorm"
)

type DB struct {
  suppress.ISuppress
}

func (s *DB) FindAuthModelByKey(db *gorm.DB, model *models.AliAuthModel) bool {
  e := db.Where("Key = ?", model.Key).First(&model)
  if e.Error == gorm.ErrRecordNotFound {
    return false
  } else if e.Error != nil {
    s.Suppress(e.Error)
  }

  return true
}
