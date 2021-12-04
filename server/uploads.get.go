package server

import (
	"github.com/Myriad-Dreamin/aliali/model"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"gorm.io/gorm"
	"time"
)

type GetUploadsRequest struct {
	Page     int64 `url:"page"`
	PageSize int64 `url:"page_size"`
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
	Items []UploadDTO `json:"items"`
}

func (srv *Server) GetUploadList(ctx *context.Context) {
	var req GetUploadsRequest
	if err := ctx.ReadQuery(&req); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		srv.Logger.Println(err.Error())
		_, _ = ctx.JSON(MessageResponse{
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

	offset := req.Page*req.PageSize - req.Page
	limit := req.PageSize

	var list []model.UploadModel
	srv.DB.Model(&model.UploadModel{}).
		Offset(int(offset)).
		Limit(int(limit)).
		Find(&list)

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
		Items: dtoList,
	})
}
