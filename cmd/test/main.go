package main

const (
	DefaultChunkSize = 1048576
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
}
