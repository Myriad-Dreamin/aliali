package ali_drive

import (
	"fmt"
	"strings"
)

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
	ProofHash     string        `json:"proof_hash"`
	PreHash       string        `json:"pre_hash"`
}

func (y *Ali) checkFileIntegrity(req *UploadFileRequest) bool {
	{
		var subReq = &ApiFileGetByIdRequest{
			DriveDirentID: req.Session.DriveDirentID,
		}
		resp := y.FileGetById(subReq)
		if resp == nil {
			return false
		}
		if resp.ContentHash != strings.ToUpper(req.Session.Hash) || resp.ContentHashName != "sha1" {
			return false
		}
		fmt.Printf("[AliYunDrive] 文件已校验: %s(%s)\n", req.Name, req.Session.Hash)
	}

	return true
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
			subReq.ProofHash = req.Session.ProofHash
			subReq.ProofHashName = "v1"
		}

		resp := y.FileCreateWithFolders(subReq)
		if resp == nil {
			return false
		}
		req.Session.DriveDirentID = resp.DriveDirentID
		req.Session.UploadID = resp.UploadID
		req.Session.PartInfoList = resp.PartInfoList
		fmt.Printf("[AliYunDrive] 文件上传开始: %s(%s) => %s\n", req.Name, req.Session.Hash, req.Session.UploadID)

		if strings.HasPrefix(req.Session.UploadID, "rapid-") {
			return y.checkFileIntegrity(req)
		}
	}

	// upload parts
	{
		var subReq = &ApiFileUploadPartRequest{}
		for i := range req.Session.PartInfoList {
			subReq.Uri = req.Session.PartInfoList[i].UploadURL
			subReq.Reader = req.File
			if y.FileUploadPart(subReq) == nil {
				return false
			}
		}
	}

	// send complete request
	{
		var subReq = &ApiFileUploadCompleteRequest{
			DriveDirentID: req.Session.DriveDirentID,
			UploadID:      req.Session.UploadID,
		}
		if y.FileUploadComplete(subReq) == nil {
			return false
		}
	}

	return y.checkFileIntegrity(req)
}
