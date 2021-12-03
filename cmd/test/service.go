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
	Session() *ali_drive.UploadSession
	FileName() string

	Size() int64
	ChunkHint() int64
	ReadAt(pos, maxLen int64) io.Reader

	Done()
}

type IUploadAliView interface {
	UploadFile(req *ali_drive.UploadFileRequest) bool
}

type Service struct {
}

const (
	UploadOK int = iota
	UploadCancelled
)

func (svc *Service) Upload(context context.Context, ali IUploadAliView, req IUploadRequest) *UploadResponse {
	var uploadingFile = ali_drive.SizedReader{
		Reader: req.ReadAt(0, req.Size()),
		Size:   req.Size(),
	}

	var ses = req.Session()
	if !ali.UploadFile(&ali_drive.UploadFileRequest{
		DriveID: ses.DriveDirentID.DriveID,
		Name:    req.FileName(),
		File:    uploadingFile,
		Session: ses,
	}) {
		return &UploadResponse{Code: UploadCancelled}
	}

	return &UploadResponse{Code: UploadOK}
}

type MockService struct {
}

func (svc *MockService) Upload(context context.Context, ali IUploadAliView, req IUploadRequest) *UploadResponse {
	return &UploadResponse{Code: 0}
}
