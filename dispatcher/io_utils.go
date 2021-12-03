package dispatcher

import (
	"bytes"
	ali_drive "github.com/Myriad-Dreamin/aliali/pkg/ali-drive"
	"github.com/Myriad-Dreamin/aliali/pkg/suppress"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

type RandReadCloser interface {
	io.ReaderAt
	io.Closer
}

type rrr struct {
	b []byte
}

func (r *rrr) Close() error {
	return nil
}

func (r *rrr) ReadAt(p []byte, off int64) (n int, err error) {
	copy(p, r.b[off:])
	if len(r.b) < len(p) {
		n = len(r.b)
		return
	}

	n = len(p)
	return n, nil
}

func NewBytesRandReader(b []byte) RandReadCloser {
	return &rrr{b}
}

type BaseUploadRequest struct {
	XSession   *ali_drive.UploadSession
	XFileName  string
	XSize      int64
	XChunkHint int64
}

func (b *BaseUploadRequest) Session() *ali_drive.UploadSession {
	return b.XSession
}

func (b *BaseUploadRequest) FileName() string {
	return b.XFileName
}

func (b *BaseUploadRequest) Size() int64 {
	return b.XSize
}

func (b *BaseUploadRequest) ChunkHint() int64 {
	return b.XChunkHint
}

type RandReaderUploadRequest struct {
	BaseUploadRequest
	r RandReadCloser
	s suppress.ISuppress
}

func (f *RandReaderUploadRequest) Done() {
	f.s.Suppress(f.r.Close())
}

func (f *RandReaderUploadRequest) ReadAt(pos, maxLen int64) io.Reader {
	sz := maxLen
	if pos+sz > f.Size() {
		sz = f.Size() - pos
	}

	var bufRaw = make([]byte, sz)
	x, err := f.r.ReadAt(bufRaw, pos)
	if err != nil {
		f.s.Suppress(err)
		return nil
	}

	if int64(x) < sz {
		sz = int64(x)
	}

	var buf = bytes.NewBuffer(bufRaw)
	return buf
}

type osFs struct {
	s string
	fs.FS
}

func (o osFs) Remove(s string) error {
	return os.Remove(filepath.Join(o.s, s))
}

func realFs() FsClearInterface {
	return osFs{"/", os.DirFS("")}
}

type ignoreFs struct {
	fs.FS
}

func (o ignoreFs) Remove(s string) error {
	return nil
}

func ignoreRemoveFs(f fs.FS) FsClearInterface {
	return ignoreFs{FS: f}
}
