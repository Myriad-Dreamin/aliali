package main

import (
	ali_notifier "github.com/Myriad-Dreamin/aliali/pkg/ali-notifier"
	"github.com/Myriad-Dreamin/aliali/pkg/suppress"
	"testing"
	"testing/fstest"
	"time"
)

type MapRmRs struct {
	fstest.MapFS
}

func (m *MapRmRs) Remove(s string) error {
	return nil
}

func TestUpload(t *testing.T) {
	var worker = NewWorker(
		MockDB(),
		WithServiceReplicate(&MockService{}),
		WithConfig(&ali_notifier.Config{
			Version: "aliyunpan/v1beta",
			AliDrive: ali_notifier.AliDriveConfig{
				RefreshToken: "123",
				DriveId:      "456",
				RootPath:     "/",
				ChunkSize:    "123456",
			},
		}))

	content := []byte("123")
	var x = NewBytesRandReader(content)

	err := worker.serveUploadRequest(&MapRmRs{map[string]*fstest.MapFile{
		"test": {
			Mode:    0644,
			ModTime: time.Time{},
		},
	}}, &ali_notifier.FsUploadRequest{
		TransactionID: 1,
		LocalPath:     "test",
		RemotePath:    "remove/test",
	}, &RandReaderUploadRequest{
		BaseUploadRequest: BaseUploadRequest{
			XDriverID:  "456",
			XFileName:  "test",
			XSize:      int64(len(content)),
			XChunkHint: 1024,
		},
		r: x,
		s: &suppress.PanicAll{},
	})
	if err != nil {
		return
	}
}
