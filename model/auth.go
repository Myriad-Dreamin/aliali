package model

import (
	"encoding/json"
	ali_drive "github.com/Myriad-Dreamin/aliali/pkg/ali-drive"
	"github.com/Myriad-Dreamin/aliali/pkg/suppress"
	"gorm.io/gorm"
)

type AliAuthModel struct {
	gorm.Model
	Key          string `gorm:"unique"`
	ExpiresLocal int64
	Raw          []byte

	memory *ali_drive.ApiAcquireTokenResponse
}

func (model *AliAuthModel) Get(s suppress.ISuppress) *ali_drive.ApiAcquireTokenResponse {
	if model.memory != nil {
		return model.memory
	}

	e := json.Unmarshal(model.Raw, &model.memory)
	if e != nil {
		s.Suppress(e)
	}

	return model.memory
}
