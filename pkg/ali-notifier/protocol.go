package ali_notifier

type FsUploadRequest struct {
	TransactionID uint64

	DriveID    string
	RemotePath string
	LocalPath  string
}
