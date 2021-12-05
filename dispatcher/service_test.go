package dispatcher

import (
	"context"
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

func newMockDispatcher(notifier *ali_notifier.Notifier, svc service.IService) *Dispatcher {
	return NewDispatcher(
		MockDB(),
		WithNotifier(notifier),
		WithServiceReplicate(svc),
		WithConfig(&ali_notifier.Config{
			Version: "aliyunpan/v1beta",
			AliDrive: ali_notifier.AliDriveConfig{
				RefreshToken: "123",
				DriveId:      "456",
				RootPath:     "/",
				ChunkSize:    "123456",
			},
		}))
}

func NewMemoryRequest(
	fsReq *ali_notifier.FsUploadRequest, contentS string) (FsClearInterface, *ali_notifier.FsUploadRequest, service.IUploadRequest) {
	content := []byte(contentS)
	var x = NewBytesRandReader(content)
	return &MapRmRs{map[string]*fstest.MapFile{
			fsReq.LocalPath: {
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
				XFileName:  fsReq.LocalPath,
				XSize:      int64(len(content)),
				XChunkHint: 1024,
			},
			R: x,
			S: &suppress.PanicAll{},
		}
}

func ExpectTaskStatus(t *testing.T, d *Dispatcher, fsReq *ali_notifier.FsUploadRequest, status int) {
	t.Helper()
	var m = &model.UploadModel{
		DriveID:    fsReq.DriveID,
		RemotePath: fsReq.RemotePath,
		LocalPath:  fsReq.LocalPath,
	}
	if !d.xdb.FindUploadRequest(d.db, m) {
		t.Error(errors.New("req not found"))
	}

	if m.Status != status {
		t.Error(errors.New("not clear"))
	}
}

func TestUpload(t *testing.T) {
	var notifier = &ali_notifier.Notifier{}

	var dispatcher = newMockDispatcher(
		notifier, &service.MockService{})

	var fsReq = &ali_notifier.FsUploadRequest{
		TransactionID: 1,
		LocalPath:     "test",
		RemotePath:    "remove/test",
	}

	notifier.Emit(fsReq)

	err := dispatcher.serveUploadRequest(NewMemoryRequest(fsReq, "123"))
	if err != nil {
		t.Error(err)
		return
	}

	// drain work
	svc := <-dispatcher.serviceQueue
	dispatcher.serviceQueue <- svc

	ExpectTaskStatus(t, dispatcher, fsReq, model.UploadStatusSettledClear)
}

func TestUploadTokenExpired(t *testing.T) {
	var notifier = &ali_notifier.Notifier{}

	var testStatus = 0

	var dispatcher = newMockDispatcher(
		notifier, &service.MockService{
			UploadHandler: func(context context.Context,
				ali service.IUploadAliView, req service.IUploadRequest) *service.UploadResponse {
				switch testStatus {
				case 0:
					if x, ok := ali.(*ali_drive.Ali); ok {
						x.InjectSemaError(401, &ali_drive.ApiErrorResponse{
							Code:      "AccessTokenInvalid",
							Message:   "injected",
							RequestId: "0000-0000-000000000000-0000",
						})
					}
					return &service.UploadResponse{
						Code: service.UploadCancelled,
					}
				case 1:
					return &service.UploadResponse{
						Code: 0,
					}
				default:
					return &service.UploadResponse{
						Code: service.UploadCancelled,
					}
				}
			},
		})

	var fsReq = &ali_notifier.FsUploadRequest{
		TransactionID: 1,
		LocalPath:     "test",
		RemotePath:    "remove/test",
	}

	notifier.Emit(fsReq)

	var serveUploadRequest = func() {
		err := dispatcher.serveUploadRequest(NewMemoryRequest(fsReq, "123"))
		if err != nil {
			t.Error(err)
			return
		}

		// drain work
		svc := <-dispatcher.serviceQueue
		dispatcher.serviceQueue <- svc
	}

	serveUploadRequest()
	//
	if _, ok := <-dispatcher.tokenInvalid; !ok {
		t.Error("token invalid not emitted")
	}

	testStatus = 1
	serveUploadRequest()

	ExpectTaskStatus(t, dispatcher, fsReq, model.UploadStatusSettledClear)
}
