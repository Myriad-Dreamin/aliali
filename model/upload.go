package model

import (
	"encoding/json"
	ali_drive "github.com/Myriad-Dreamin/aliali/pkg/ali-drive"
	"github.com/Myriad-Dreamin/aliali/pkg/suppress"
	"gorm.io/gorm"
	"time"
)

const (
	UploadStatusUninitialized   int = 0
	UploadStatusInitialized         = 1
	UploadStatusUploading           = 2
	UploadStatusUploaded            = 3
	UploadStatusSetting             = 4
	UploadStatusSettledClear        = 5
	UploadStatusSettledNotClear     = 6

	UploadStatusSettledExit = 1000
)

type UploadModel struct {
	ID         uint64 `gorm:"primarykey"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	Status     int            `gorm:"column:status"`
	DriveID    string
	RemotePath string
	LocalPath  string
	Raw        []byte

	uploadData *ali_drive.UploadSession
}

func (model *UploadModel) Get(s suppress.ISuppress) *ali_drive.UploadSession {
	if model.uploadData != nil {
		return model.uploadData
	}

	e := json.Unmarshal(model.Raw, &model.uploadData)
	if e != nil {
		s.Suppress(e)
		return nil
	}

	return model.uploadData
}
