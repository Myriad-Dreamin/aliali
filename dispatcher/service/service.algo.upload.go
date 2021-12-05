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

	hashBytes, err = computeHash(sha1.New(), RangeReader(req, req.ChunkHint(), req.Size()))
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
