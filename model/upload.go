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
	UploadStatusSettledClear        = 4
	UploadStatusSettledNotClear     = 5

	UploadStatusSettledExit            = 1000
	UploadStatusSettledExitFileFlyAway = 1001
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
	Group      string
	Raw        []byte `gorm:"column:raw"`

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

func (model *UploadModel) Set(s suppress.ISuppress, session *ali_drive.UploadSession) bool {
	b, e := json.Marshal(&session)
	if e != nil {
		s.Suppress(e)
		return false
	}

	model.Raw = b
	model.uploadData = session
	return true
}
