package server

import (
	"github.com/Myriad-Dreamin/aliali/model"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"gorm.io/gorm"
	"strings"
	"time"
)

type GetUploadsRequest struct {
	Deleted      *bool  `url:"deleted"`
	GroupPattern string `url:"group"`
	Page         int64  `url:"page"`
	PageSize     int64  `url:"page_size"`
}

type UploadDTO struct {
	ID         uint64         `json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at"`
	Status     int            `json:"status"`
	DriveID    string         `json:"drive_id"`
	RemotePath string         `json:"remote_path"`
	LocalPath  string         `json:"local_path"`
	UploadID   string         `json:"upload_id"`
	Hash       string         `json:"hash"`
	PreHash    string         `json:"pre_hash"`
}

type GetUploadsResponse struct {
	Items   []UploadDTO `json:"items"`
	Count   int64       `json:"count"`
	Current int64       `json:"current"`
}

func (srv *Server) GetUploadList(ctx *context.Context) {
	var req GetUploadsRequest
	if err := ctx.ReadQuery(&req); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		srv.Logger.Println(err.Error())
		_, _ = ctx.JSON(StdResponse{
			Code:    CodeInvalidParams,
			Message: "解析请求参数发生错误，查看后台日志了解内容...",
		})
		return
	}

	if req.Page == 0 {
		req.Page = 1
	}

	if req.PageSize == 0 {
		req.PageSize = 20
	}

	offset := req.Page*req.PageSize - req.PageSize
	limit := req.PageSize

	var list []model.UploadModel
	db := srv.DB.Model(&model.UploadModel{}).Debug().
		Offset(int(offset)).
		Limit(int(limit))
	if strings.Contains(req.GroupPattern, "%") {
		db = db.Where("group like ?", req.GroupPattern)
	} else if len(req.GroupPattern) != 0 {
		db = db.Where("group = ?", req.GroupPattern)
	}

	if req.Deleted != nil && *req.Deleted {
		db = db.Unscoped().Where("deleted_at is not null")
	}

	db.Find(&list)

	var cnt int64
	db.Count(&cnt)

	var dtoList []UploadDTO
	for i := range list {
		var x = UploadDTO{
			ID:         list[i].ID,
			CreatedAt:  list[i].CreatedAt,
			UpdatedAt:  list[i].UpdatedAt,
			DeletedAt:  list[i].DeletedAt,
			Status:     list[i].Status,
			DriveID:    list[i].DriveID,
			RemotePath: list[i].RemotePath,
			LocalPath:  list[i].LocalPath,
		}
		if len(list[i].Raw) != 0 {
			ses := list[i].Get(srv.S)
			x.UploadID = ses.UploadID
			x.Hash = ses.Hash
			x.PreHash = ses.PreHash
		}
		dtoList = append(dtoList, x)
	}

	_, _ = ctx.JSON(&GetUploadsResponse{
		Items:   dtoList,
		Count:   cnt,
		Current: offset,
	})
}
