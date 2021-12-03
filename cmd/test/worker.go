package main

import (
	"github.com/Myriad-Dreamin/aliali/model"
	ali_drive "github.com/Myriad-Dreamin/aliali/pkg/ali-drive"
	ali_notifier "github.com/Myriad-Dreamin/aliali/pkg/ali-notifier"
	"github.com/Myriad-Dreamin/aliali/pkg/suppress"
	"gorm.io/gorm"
)

type Worker struct {
	cfg         *ali_notifier.Config
	auth        *model.AliAuthModel
	httpHeaders [][2]string

	ali       *ali_drive.Ali
	authedAli *ali_drive.Ali
	s         suppress.ISuppress
	db        *gorm.DB

	fileUploads  chan *FsUploadRequest
	serviceQueue chan IService
}

type Option = func(w *Worker) *Worker

func MockDB() Option {
	return func(w *Worker) *Worker {
		w.db = w.openMock()
		if w.db == nil {
			return nil
		}
		return w
	}
}

func WithConfig(cfg *ali_notifier.Config) Option {
	return func(w *Worker) *Worker {
		w.cfg = cfg
		return w
	}
}

func WithServiceReplicate(services ...IService) Option {
	return func(w *Worker) *Worker {
		w.serviceQueue = make(chan IService, len(services))
		for i := range services {
			w.serviceQueue <- services[i]
		}
		return w
	}
}

func NewWorker(options ...Option) *Worker {
	s := suppress.PanicAll{}

	var httpHeaders [][2]string
	httpHeaders = append(httpHeaders, [2]string{"origin", "https://aliyundrive.com"})
	httpHeaders = append(httpHeaders, [2]string{"referer", "https://aliyundrive.com"})

	var w = &Worker{
		s:           s,
		fileUploads: make(chan *FsUploadRequest, 10),
		httpHeaders: httpHeaders,
	}

	for _, option := range options {
		w = option(w)
		if w == nil {
			return nil
		}
	}

	w.ali = w.makeAliClient()
	if w.cfg == nil {
		w.cfg = w.syncConfig()
	}
	if w.db == nil {
		w.db = w.openDB()
	}
	if w.serviceQueue == nil {
		w.serviceQueue = make(chan IService, 1)
		w.serviceQueue <- &Service{}
	}

	if w.ali == nil || w.db == nil || w.cfg == nil || w.serviceQueue == nil {
		return nil
	}

	return w
}

func (w *Worker) warnOnce(err error) {
	w.s.Suppress(err)
}
