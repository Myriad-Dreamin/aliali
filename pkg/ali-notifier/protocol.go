package ali_notifier

type FsUploadRequest struct {
	TransactionID uint64

	Group      string
	DriveID    string
	RemotePath string
	LocalPath  string
}
