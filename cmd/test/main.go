package main

import (
	ali_notifier "github.com/Myriad-Dreamin/aliali/pkg/ali-notifier"
	"github.com/Myriad-Dreamin/aliali/pkg/suppress"
)

const (
	DefaultChunkSize  = 1048576
	DefaultConfigPath = "config.yaml"
)

func (w *Worker) loop() {
	for {
		select {
		case req := <-w.fileUploads:
			if w.authExpired() {
				w.refreshAuth()
			}
			if err := w.serveFsUploadRequest(req); err != nil {
				w.s.Suppress(err)
			}
		}
	}
}

func main() {
	//var mockFile = bytes.NewBuffer(nil)
	//mockFile.WriteString("test2")

	var s = suppress.PanicAll{}

	notifier := ali_notifier.BiliRecorderNotifier{}
	s.Suppress(notifier.Run())
}
