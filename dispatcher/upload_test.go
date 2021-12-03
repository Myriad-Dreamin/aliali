package dispatcher

import (
	"errors"
	"github.com/Myriad-Dreamin/aliali/dispatcher/service"
	"github.com/Myriad-Dreamin/aliali/model"
	ali_drive "github.com/Myriad-Dreamin/aliali/pkg/ali-drive"
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
	var notifier = &ali_notifier.Notifier{}

	var dispatcher = NewDispatcher(
		MockDB(),
		WithNotifier(notifier),
		WithServiceReplicate(&service.MockService{}),
		WithConfig(&ali_notifier.Config{
			Version: "aliyunpan/v1beta",
			AliDrive: ali_notifier.AliDriveConfig{
				RefreshToken: "123",
				DriveId:      "456",
				RootPath:     "/",
				ChunkSize:    "123456",
			},
		}))

	var fsReq = &ali_notifier.FsUploadRequest{
		TransactionID: 1,
		LocalPath:     "test",
		RemotePath:    "remove/test",
	}

	notifier.Emit(fsReq)

	content := []byte("123")
	var x = NewBytesRandReader(content)

	err := dispatcher.serveUploadRequest(&MapRmRs{map[string]*fstest.MapFile{
		"test": {
			Mode:    0644,
			ModTime: time.Time{},
		},
	}}, fsReq, &RandReaderUploadRequest{
		BaseUploadRequest: BaseUploadRequest{
			XSession: &ali_drive.UploadSession{
				DriveDirentID: ali_drive.DriveDirentID{
					DriveID: "456",
				},
			},
			XFileName:  "test",
			XSize:      int64(len(content)),
			XChunkHint: 1024,
		},
		R: x,
		S: &suppress.PanicAll{},
	})
	if err != nil {
		t.Error(err)
		return
	}

	// drain work
	svc := <-dispatcher.serviceQueue
	dispatcher.serviceQueue <- svc

	var m = &model.UploadModel{
		DriveID:    fsReq.DriveID,
		RemotePath: fsReq.RemotePath,
		LocalPath:  fsReq.LocalPath,
	}
	if !dispatcher.xdb.FindUploadRequest(dispatcher.db, m) {
		t.Error(errors.New("req not found"))
	}

	if m.Status != model.UploadStatusSettledClear {
		t.Error(errors.New("not clear"))
	}
}