package dispatcher

import (
	"github.com/Myriad-Dreamin/aliali/model"
	ali_notifier "github.com/Myriad-Dreamin/aliali/pkg/ali-notifier"
)

func fsUploadReq2Model(req *ali_notifier.FsUploadRequest) *model.UploadModel {
	return &model.UploadModel{
		ID:         req.TransactionID,
		Group:      req.Group,
		DriveID:    req.DriveID,
		RemotePath: req.RemotePath,
		LocalPath:  req.LocalPath,
	}
}

func uploadModel2fsReq(req *model.UploadModel) *ali_notifier.FsUploadRequest {
	return &ali_notifier.FsUploadRequest{
		TransactionID: req.ID,
		Group:         req.Group,
		DriveID:       req.DriveID,
		RemotePath:    req.RemotePath,
		LocalPath:     req.LocalPath,
	}
}
