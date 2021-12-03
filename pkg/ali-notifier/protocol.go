package ali_notifier

type FsUploadRequest struct {
	TransactionID uint64

	LocalPath  string
	RemotePath string
}
