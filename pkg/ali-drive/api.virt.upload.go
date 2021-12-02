package ali_drive

type UploadFileRequest struct {
	DriveID string      `json:"drive_id"`
	Name    string      `json:"name"`
	File    SizedReader `json:"-"`
}

func (y *Ali) UploadFile(req *UploadFileRequest) {
	var dirent DriveDirentID
	var parts = []PartInfo{
		{PartNumber: 1},
	}
	var uploadId string

	// create file uploading proof
	{
		var subReq = &ApiFileCreateWithFoldersRequest{
			DriveID:       req.DriveID,
			ParentFileID:  "root",
			Name:          req.Name,
			CheckNameMode: "auto_rename",
			Type:          "file",
			PartInfoList:  parts,
			Size:          req.File.Size,
		}
		resp := y.FileCreateWithFolders(subReq)
		dirent = resp.DriveDirentID
		uploadId = resp.UploadID
		parts = resp.PartInfoList
	}

	// upload parts
	{
		var subReq = &ApiFileUploadPartRequest{}
		for i := range parts {
			subReq.Uri = parts[i].UploadURL
			subReq.Reader = req.File
			y.FileUploadPart(subReq)
		}
	}

	// send complete request
	{
		var subReq = &ApiFileUploadCompleteRequest{
			DriveDirentID: dirent,
			UploadID:      uploadId,
		}
		y.FileUploadComplete(subReq)
	}
}
