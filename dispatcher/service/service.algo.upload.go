package service

import (
	"context"
	ali_drive "github.com/Myriad-Dreamin/aliali/pkg/ali-drive"
	"io"
	"log"
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

type UploadImpl struct {
	Logger *log.Logger
}

const (
	UploadOK int = iota
	UploadCancelled
)

func (svc *UploadImpl) Upload(context context.Context, ali IUploadAliView, req IUploadRequest) *UploadResponse {
	var uploadingFile = ali_drive.SizedReader{
		Reader: req.ReadAt(0, req.Size()),
		Size:   req.Size(),
	}

	var ses = req.Session()
	svc.Logger.Printf("begin file upload session: %v", req)
	if !ali.UploadFile(&ali_drive.UploadFileRequest{
		DriveID: ses.DriveDirentID.DriveID,
		Name:    req.FileName(),
		File:    uploadingFile,
		Session: ses,
	}) {
		svc.Logger.Printf("file upload session cancelled: %v", req)
		return &UploadResponse{Code: UploadCancelled}
	}

	svc.Logger.Printf("successful file upload session: %v", req)
	return &UploadResponse{Code: UploadOK}
}

type MockService struct {
}

func (svc *MockService) Upload(context context.Context, ali IUploadAliView, req IUploadRequest) *UploadResponse {
	return &UploadResponse{Code: 0}
}
