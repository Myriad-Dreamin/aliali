package dispatcher

import (
	"bytes"
	ali_drive "github.com/Myriad-Dreamin/aliali/pkg/ali-drive"
	ali_utils "github.com/Myriad-Dreamin/aliali/pkg/ali-utils"
	"github.com/Myriad-Dreamin/aliali/pkg/suppress"
	"io"
	"io/fs"
)

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
	R ali_utils.RandReadCloser
	S suppress.ISuppress
}

func (f *RandReaderUploadRequest) Done() {
	f.S.Suppress(f.R.Close())
}

func (f *RandReaderUploadRequest) ReadAt(pos, maxLen int64) io.Reader {
	sz := maxLen
	if pos+sz > f.Size() {
		sz = f.Size() - pos
	}

	var bufRaw = make([]byte, sz)
	x, err := f.R.ReadAt(bufRaw, pos)
	if err != nil {
		f.S.Suppress(err)
		return nil
	}

	if int64(x) < sz {
		sz = int64(x)
	}

	var buf = bytes.NewBuffer(bufRaw)
	return buf
}

type ignoreFs struct {
	fs.FS
}

func (o ignoreFs) Remove(s string) error {
	return nil
}

func ignoreRemoveFs(f fs.FS) ali_utils.FsClearInterface {
	return ignoreFs{FS: f}
}
