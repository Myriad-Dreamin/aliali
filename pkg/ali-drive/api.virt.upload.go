package ali_drive

import "fmt"

type UploadFileRequest struct {
	DriveID string         `json:"drive_id"`
	Name    string         `json:"name"`
	File    SizedReader    `json:"-"`
	Session *UploadSession `json:"session"`
}

type UploadSession struct {
	DriveDirentID DriveDirentID `json:"dirent_id"`
	PartInfoList  []PartInfo    `json:"part_info_list"`
	UploadID      string        `json:"upload_id"`
	Hash          string        `json:"hash"`
	PreHash       string        `json:"pre_hash"`
}

func (y *Ali) UploadFile(req *UploadFileRequest) bool {
	if req.Session == nil {
		req.Session = &UploadSession{}
	}

	if req.Session.PartInfoList == nil {
		req.Session.PartInfoList = []PartInfo{{PartNumber: 1}}
	}

	// create file uploading proof
	{
		var subReq = &ApiFileCreateWithFoldersRequest{
			DriveID:       req.DriveID,
			ParentFileID:  "root",
			Name:          req.Name,
			CheckNameMode: "auto_rename",
			Type:          "file",
			PartInfoList:  req.Session.PartInfoList,
			Size:          req.File.Size,
		}

		if len(req.Session.Hash) != 0 {
			subReq.ContentHash = req.Session.Hash
			subReq.ContentHashName = "sha1"
		}

		fmt.Println("hash?", req.Session.Hash, req.Session.PreHash, subReq)
		resp := y.FileCreateWithFolders(subReq)
		req.Session.DriveDirentID = resp.DriveDirentID
		req.Session.UploadID = resp.UploadID
		req.Session.PartInfoList = resp.PartInfoList
	}

	// upload parts
	{
		var subReq = &ApiFileUploadPartRequest{}
		for i := range req.Session.PartInfoList {
			subReq.Uri = req.Session.PartInfoList[i].UploadURL
			subReq.Reader = req.File
			y.FileUploadPart(subReq)
		}
	}

	// send complete request
	{
		var subReq = &ApiFileUploadCompleteRequest{
			DriveDirentID: req.Session.DriveDirentID,
			UploadID:      req.Session.UploadID,
		}
		y.FileUploadComplete(subReq)
	}

	return true
}
