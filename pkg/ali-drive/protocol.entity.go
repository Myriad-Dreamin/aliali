package ali_drive

import "time"

type DriveDirentID struct {
	DriveID string `json:"drive_id"`
	FileID  string `json:"file_id"`
}

type PartInfo struct {
	PartNumber int    `json:"part_number"`
	UploadURL  string `json:"upload_url"`
}

type File struct {
	DriveID         string    `json:"drive_id"`
	DomainID        string    `json:"domain_id"`
	FileID          string    `json:"file_id"`
	Name            string    `json:"name"`
	Type            string    `json:"type"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Hidden          bool      `json:"hidden"`
	Starred         bool      `json:"starred"`
	Status          string    `json:"status"`
	UserMeta        string    `json:"user_meta,omitempty"`
	ParentFileID    string    `json:"parent_file_id"`
	EncryptMode     string    `json:"encrypt_mode"`
	ContentType     string    `json:"content_type,omitempty"`
	FileExtension   string    `json:"file_extension,omitempty"`
	MimeType        string    `json:"mime_type,omitempty"`
	MimeExtension   string    `json:"mime_extension,omitempty"`
	Size            int64     `json:"size,omitempty"`
	Crc64Hash       string    `json:"crc64_hash,omitempty"`
	ContentHash     string    `json:"content_hash,omitempty"`
	ContentHashName string    `json:"content_hash_name,omitempty"`
	DownloadURL     string    `json:"download_url,omitempty"`
	URL             string    `json:"url,omitempty"`
	Thumbnail       string    `json:"thumbnail,omitempty"`
	Category        string    `json:"category,omitempty"`
	PunishFlag      int       `json:"punish_flag,omitempty"`
}
