package ali_utils

import (
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
	if off > int64(len(r.b)) {
		return 0, io.EOF
	}
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

type FsClearInterface interface {
	fs.FS
	Remove(path string) error
}

type osFs struct {
	s string
	fs.FS
}

func (o osFs) Remove(s string) error {
	return os.Remove(filepath.Join(o.s, s))
}

func RealFs() FsClearInterface {
	return osFs{"/", os.DirFS("")}
}
