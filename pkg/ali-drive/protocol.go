package ali_drive

import (
	"io"
)

type ApiFileCreateWithProofRequest struct {
	DriveID       string      `json:"drive_id"`
	PartInfoList  []*PartInfo `json:"part_info_list"`
	ParentFileID  string      `json:"parent_file_id"`
	Name          string      `json:"name"`
	Type          string      `json:"type"`
	CheckNameMode string      `json:"check_name_mode"`
	Size          int64       `json:"size"`
	PreHash       string      `json:"pre_hash"`
}

type ApiFileCreateWithProofResponse struct {
	UploadID     string      `json:"upload_id"`
	FileID       string      `json:"file_id"`
	PartInfoList []*PartInfo `json:"part_info_list"`
}

type ApiFileUploadPartRequest struct {
	Uri    string    `json:"-"`
	Reader io.Reader `json:"-"`
}

type ApiFileUploadPartResponse struct {
}

type ApiFileUploadCompleteRequest struct {
	DriveDirentID

	UploadID string `json:"upload_id"`
}

type ApiFileUploadCompleteResponse struct {
}

type ApiFileSearchRequest struct {
	DriveID               string `json:"drive_id"`
	Limit                 int    `json:"limit"`
	Query                 string `json:"query"`
	ImageThumbnailProcess string `json:"image_thumbnail_process"`
	ImageURLProcess       string `json:"image_url_process"`
	VideoThumbnailProcess string `json:"video_thumbnail_process"`
	OrderBy               string `json:"order_by"`
}

type ApiFileSearchResponse struct {
	Message    string  `json:"message"`
	Items      []*File `json:"items"`
	NextMarker string  `json:"next_marker"`
}

type ApiFileGetByIdRequest struct {
	DriveDirentID
}

type ApiFileGetByIdResponse struct {
	File
	Message string `json:"message"`
	Trashed bool   `json:"trashed"`
}

type ApiFileGetPathByIdRequest struct {
	DriveDirentID
}

type ApiFileGetPathByIdResponse struct {
	Message string  `json:"message"`
	Items   []*File `json:"items"`
}

type ApiFileDeleteByIdRequest struct {
	DriveDirentID
}

type ApiFileDeleteByIdResponse struct {
	DriveDirentID

	Message     string `json:"message"`
	DomainID    string `json:"domain_id"`
	AsyncTaskID string `json:"async_task_id"`
}

type ApiFileCreateWithFoldersRequest struct {
	DriveID      string `json:"drive_id"`
	ParentFileID string `json:"parent_file_id"`

	Name          string `json:"name"`
	CheckNameMode string `json:"check_name_mode"`
	Type          string `json:"type"`

	PartInfoList    []*PartInfo `json:"part_info_list"`
	Size            int64       `json:"size"`
	ContentHash     string      `json:"content_hash"`
	ContentHashName string      `json:"content_hash_name"`
	ProofHash       string      `json:"proof_code"`
	ProofHashName   string      `json:"proof_version"`
}

type ApiFileCreateWithFoldersResponse struct {
	DriveDirentID

	Message      string `json:"message"`
	ParentFileID string `json:"parent_file_id"`
	Type         string `json:"type"`
	DomainID     string `json:"domain_id"`
	FileName     string `json:"file_name"`
	EncryptMode  string `json:"encrypt_mode"`
}
