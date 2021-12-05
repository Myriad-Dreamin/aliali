package ali_drive

import (
	"io"
)

type ApiErrorResponse struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestId string `json:"requestId"`
}

type ApiAcquireTokenRequest struct {
	GrantType    string `json:"grant_type"`
	RefreshToken string `json:"refresh_token"`
}

type ApiAcquireTokenResponse struct {
	AccessToken    string   `json:"access_token"`
	Avatar         string   `json:"avatar"`
	DefaultDriveId string   `json:"default_drive_id"`
	DeviceId       string   `json:"device_id"`
	ExistLink      []string `json:"exist_link"`
	ExpireTime     string   `json:"expire_time"`
	ExpiresIn      int      `json:"expires_in"`
	IsFirstLogin   bool     `json:"is_first_login"`
	NeedLink       bool     `json:"need_link"`
	NeedRpVerify   bool     `json:"need_rp_verify"`
	NickName       string   `json:"nick_name"`
	PinSetup       bool     `json:"pin_setup"`
	RefreshToken   string   `json:"refresh_token"`
	Role           string   `json:"role"`
	State          string   `json:"state"`
	Status         string   `json:"status"`
	TokenType      string   `json:"token_type"`
	UserId         string   `json:"user_id"`
	UserName       string   `json:"user_name"`
}

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

	PartInfoList    []PartInfo `json:"part_info_list"`
	Size            int64      `json:"size"`
	ContentHash     string     `json:"content_hash,omitempty"`
	ContentHashName string     `json:"content_hash_name,omitempty"`
	ProofHash       string     `json:"proof_code,omitempty"`
	ProofHashName   string     `json:"proof_version,omitempty"`
}

type ApiFileCreateWithFoldersResponse struct {
	DriveDirentID

	PartInfoList []PartInfo `json:"part_info_list"`
	Message      string     `json:"message"`
	ParentFileID string     `json:"parent_file_id"`
	UploadID     string     `json:"upload_id"`
	Type         string     `json:"type"`
	DomainID     string     `json:"domain_id"`
	FileName     string     `json:"file_name"`
	EncryptMode  string     `json:"encrypt_mode"`
}

type ApiFileGetListRequest struct {
	GetAll                bool   `json:"get_all"`
	DriveID               string `json:"drive_id"`
	ParentFileID          string `json:"parent_file_id"`
	Marker                string `json:"marker"`
	Limit                 int    `json:"limit"`
	All                   bool   `json:"all"`
	URLExpireSec          int    `json:"url_expire_sec"`
	ImageThumbnailProcess string `json:"image_thumbnail_process"`
	ImageURLProcess       string `json:"image_url_process"`
	VideoThumbnailProcess string `json:"video_thumbnail_process"`
	Fields                string `json:"fields"`
	OrderBy               string `json:"order_by"`
	OrderDirection        string `json:"order_direction"`
}

type ApiFileGetListResponse struct {
	Message    string  `json:"message"`
	Items      []*File `json:"items"`
	NextMarker string  `json:"next_marker"`
}
