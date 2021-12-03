package service

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
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
	UploadHashComputingFailed
	UploadCancelled
)

type XReaderAt interface {
	ReadAt(pos, maxLen int64) io.Reader
}

type rangeReaderBuf struct {
	r             XReaderAt
	chunkSize     int64
	size          int64
	pos           int64
	currentReader io.Reader
}

func (r *rangeReaderBuf) Read(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}
	if r.currentReader == nil {
		if !r.nextReader() {
			return 0, io.EOF
		}
	}

	if n, err = r.currentReader.Read(p); err != nil {
		if err == io.EOF {
			if !r.nextReader() {
				return n, err
			}

			p = p[n:]
			nn, err2 := r.Read(p)
			nn += n
			if nn != 0 && err2 == io.EOF {
				return nn, nil
			}
			return nn, err2
		}

		return
	}

	return
}

func (r *rangeReaderBuf) nextReader() bool {
	if r.pos >= r.size {
		return false
	}

	chunkSize := r.chunkSize
	if chunkSize > r.size-r.pos {
		chunkSize = r.size - r.pos
	}
	r.currentReader = r.r.ReadAt(r.pos, chunkSize)
	if r.currentReader == nil {
		return false
	}
	r.pos += chunkSize
	return true
}

func RangeReader(r XReaderAt, chunkSize, sz int64) io.Reader {
	if chunkSize == 0 {
		chunkSize = 128 * 1024
	}
	return &rangeReaderBuf{r: r, chunkSize: chunkSize, size: sz}
}

func (svc *UploadImpl) Upload(context context.Context, ali IUploadAliView, req IUploadRequest) *UploadResponse {
	var uploadingFile = ali_drive.SizedReader{
		Reader: RangeReader(req, req.ChunkHint(), req.Size()),
		Size:   req.Size(),
	}

	var ses = req.Session()
	svc.Logger.Printf("begin file upload session: %v", req)
	if len(ses.Hash) == 0 {
		err := svc.computeHash(ses, req)
		if err != nil {
			return &UploadResponse{Code: UploadHashComputingFailed, Err: err}
		}
	}

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

func computeSha1(r io.Reader) ([]byte, error) {
	h := sha1.New()
	_, err := io.Copy(h, r)
	if err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

func (svc *UploadImpl) computeHash(ses *ali_drive.UploadSession, req IUploadRequest) error {
	var hash []byte
	var err error

	hash, err = computeSha1(req.ReadAt(0, 1024))
	if err != nil {
		return err
	}
	ses.PreHash = hex.EncodeToString(hash)

	hash, err = computeSha1(RangeReader(req, req.ChunkHint(), req.Size()))
	if err != nil {
		ses.PreHash = ""
		return err
	}
	ses.Hash = hex.EncodeToString(hash)
	return nil
}

type MockService struct {
}

func (svc *MockService) Upload(context context.Context, ali IUploadAliView, req IUploadRequest) *UploadResponse {
	return &UploadResponse{Code: 0}
}
