package dispatcher

import (
	"fmt"
	"github.com/Myriad-Dreamin/aliali/model"
	ali_drive "github.com/Myriad-Dreamin/aliali/pkg/ali-drive"
	ali_notifier "github.com/Myriad-Dreamin/aliali/pkg/ali-notifier"
	"runtime"
	"time"
)

type ServiceSuppress struct {
	*Dispatcher
	*ServiceContext
}

func (s *ServiceSuppress) Suppress(err error) {
	if err != nil {
		switch e := err.(type) {
		case *ali_drive.AliSemaError:
			switch e.Code {
			case "AccessTokenInvalid":
				s.tokenInvalid <- true
			default:
				panic(err)
			}
		case *ali_drive.AliHttpError:
			s.WarnOnce(err)
		default:
			panic(err)
		}
	}
}

func (s *ServiceSuppress) WarnOnce(err error) {
	if err != nil {
		// s.Warnings = append(s.Warnings, err)
		var b = make([]byte, 1024)
		b = b[:runtime.Stack(b, false)]
		s.Dispatcher.logger.Printf("warning occurs: %s:\n%s\n", err.Error(), string(b))
	}
}

func (d *Dispatcher) waitService() *ServiceContext {
	ctx := <-d.serviceQueue

	if ctx.S == nil {
		ctx.S = &ServiceSuppress{Dispatcher: d, ServiceContext: ctx}
	}
	if ctx.authedAli == nil {
		ctx.authedAli = d.makeAliClient(ctx.S)
	}

	if d.authedAli != nil {
		ctx.authedAli.SetAccessToken(d.authedAli.GetAccessToken())
	}
	return ctx
}

func (d *Dispatcher) Loop() {
	d.logger.Printf("main dispatcher entering forever loop")
	tick := time.NewTicker(time.Second)
	var stackModel model.UploadModel

	var getReq = func() *ali_notifier.FsUploadRequest {
		return &ali_notifier.FsUploadRequest{
			TransactionID: stackModel.ID,
			DriveID:       stackModel.DriveID,
			RemotePath:    stackModel.RemotePath,
			LocalPath:     stackModel.LocalPath,
		}
	}

	for {
		if !d.xdb.FindMatchedStatusRequest(d.db, &stackModel, model.UploadStatusUploading) {
			break
		}
		fmt.Println(getReq())
		d.xdb.TransitUploadStatus(d.db, getReq(), model.UploadStatusUploading, model.UploadStatusInitialized)
	}

	for {
		if !d.xdb.FindMatchedStatusRequest(d.db, &stackModel, model.UploadStatusUploaded) {
			break
		}
		fmt.Println(getReq())
		d.xdb.TransitUploadStatus(d.db, getReq(), model.UploadStatusUploaded, model.UploadStatusInitialized)
	}

	for {
		select {
		case req := <-d.fileUploads:
			d.logger.Printf("receive file upload request: %v", req)
			if d.authExpired() {
				d.refreshAuth()
			}
			if err := d.serveFsUploadRequest(req); err != nil {
				d.s.Suppress(err)
			}
		case <-d.tokenInvalid:
			if d.authExpired() {
				d.refreshAuth()
			}
		case <-tick.C:
			if d.xdb.FindMatchedStatusRequest(d.db, &stackModel, model.UploadStatusInitialized) {
				d.fileUploads <- getReq()
			}
		}
	}
}
