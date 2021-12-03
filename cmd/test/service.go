package main

import (
	"context"
	ali_drive "github.com/Myriad-Dreamin/aliali/pkg/ali-drive"
	"io"
)

type IService interface {
	Upload(context context.Context, ali IUploadAliView, req IUploadRequest) *UploadResponse
}

type UploadResponse struct {
	Code int
	Err  error
}

type IUploadRequest interface {
	DriverID() string
	FileName() string

	Size() int64
	ChunkHint() int64
	ReadAt(pos, maxLen int64) io.Reader

	Done()
}

type IUploadAliView interface {
	UploadFile(req *ali_drive.UploadFileRequest)
}

type Service struct {
}

func (svc *Service) Upload(context context.Context, ali IUploadAliView, req IUploadRequest) *UploadResponse {
	var uploadingFile = ali_drive.SizedReader{
		Reader: req.ReadAt(0, req.Size()),
		Size:   req.Size(),
	}

	ali.UploadFile(&ali_drive.UploadFileRequest{
		DriveID: req.DriverID(),
		Name:    req.FileName(),
		File:    uploadingFile,
	})

	return &UploadResponse{Code: 0}
}

type MockService struct {
}

func (svc *MockService) Upload(context context.Context, ali IUploadAliView, req IUploadRequest) *UploadResponse {
	return &UploadResponse{Code: 0}
}
