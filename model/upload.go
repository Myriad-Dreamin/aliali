package model

import (
	"encoding/json"
	ali_drive "github.com/Myriad-Dreamin/aliali/pkg/ali-drive"
	"github.com/Myriad-Dreamin/aliali/pkg/suppress"
	"gorm.io/gorm"
)

const (
	UploadStatusUninitialized   int = 0
	UploadStatusUploaded            = 1
	UploadStatusSettledClear        = 2
	UploadStatusSettledNotClear     = 3

	UploadStatusSettledExit = 1000
)

type UploadModel struct {
	gorm.Model
	Status     int
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
