package service

import (
	"bytes"
	"context"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	ali_drive "github.com/Myriad-Dreamin/aliali/pkg/ali-drive"
	"github.com/Myriad-Dreamin/aliali/pkg/ali-utils"
	"hash"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"strings"
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
	GetAccessToken() string
}

type UploadImpl struct {
	Logger *log.Logger
}

const (
	UploadOK int = iota
	UploadHashComputingFailed
	UploadCancelled
)

func (svc *UploadImpl) Upload(context context.Context, ali IUploadAliView, req IUploadRequest) *UploadResponse {
	var uploadingFile = ali_drive.SizedReader{
		Reader: ali_utils.NewRangeReader(req, req.ChunkHint(), req.Size()),
		Size:   req.Size(),
	}

	var ses = req.Session()
	svc.Logger.Printf("begin file upload session: %v", req)
	if len(ses.Hash) == 0 {
		err := svc.computeHash(ali.GetAccessToken(), ses, req)
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

func computeHash(hash hash.Hash, r io.Reader) ([]byte, error) {
	_, err := io.Copy(hash, r)
	if err != nil {
		return nil, err
	}
	return hash.Sum(nil), nil
}

func (svc *UploadImpl) computeHash(accessToken string, ses *ali_drive.UploadSession, req IUploadRequest) error {
	var hashBytes []byte
	var err error
	accessToken = strings.TrimPrefix(accessToken, "Bearer ")

	hashBytes, err = computeHash(sha1.New(), req.ReadAt(0, 1024))
	if err != nil {
		return err
	}
	ses.PreHash = hex.EncodeToString(hashBytes)

	hashBytes, err = computeHash(sha1.New(), ali_utils.NewRangeReader(req, req.ChunkHint(), req.Size()))
	if err != nil {
		ses.PreHash = ""
		return err
	}
	ses.Hash = hex.EncodeToString(hashBytes)

	if req.Size() != 0 {
		m := url.PathEscape(url.PathEscape(accessToken))
		hashBytes, err = computeHash(md5.New(), bytes.NewReader([]byte(m)))
		if err != nil {
			ses.PreHash = ""
			ses.Hash = ""
			return err
		}

		start := binary.BigEndian.Uint64(hashBytes) % uint64(req.Size())
		hashBytes, err = ioutil.ReadAll(req.ReadAt(int64(start), 8))
		if err != nil {
			ses.PreHash = ""
			ses.Hash = ""
			return err
		}

		ses.ProofHash = base64.URLEncoding.EncodeToString(hashBytes)
	}

	return nil
}

type MockService struct {
	UploadHandler func(context context.Context, ali IUploadAliView, req IUploadRequest) *UploadResponse
}

func (svc *MockService) Upload(context context.Context, ali IUploadAliView, req IUploadRequest) *UploadResponse {
	if svc.UploadHandler != nil {
		return svc.UploadHandler(context, ali, req)
	}
	return &UploadResponse{Code: 0}
}
