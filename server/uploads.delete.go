package server

import (
	"github.com/Myriad-Dreamin/aliali/model"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

type DeleteUploadRequest struct {
	ID uint64 `json:"id"`
}

func (srv *Server) DeleteUpload(ctx *context.Context) {
	var req DeleteUploadRequest
	if err := ctx.ReadQuery(&req); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		srv.Logger.Println(err.Error())
		_, _ = ctx.JSON(StdResponse{
			Code:    CodeInvalidParams,
			Message: "解析请求参数发生错误，查看后台日志了解内容...",
		})
		return
	}

	if req.ID == 0 {
		_, _ = ctx.JSON(StdResponse{
			Code:    CodeNoSuchId,
			Message: "请求删除的上传项ID是空的",
		})
		return
	}

	db := srv.DB.Model(&model.UploadModel{
		ID: req.ID,
	}).Delete(nil)

	if db.Error != nil {
		srv.Logger.Println(db.Error.Error())
		_, _ = ctx.JSON(StdResponse{
			Code:    CodeDBExecutionError,
			Message: "删除失败，查看后台日志了解内容...",
		})
	}

	if db.RowsAffected == 0 {
		_, _ = ctx.JSON(&StdResponse{
			Code:    CodeOK,
			Message: "可能已经删过了",
		})
		return
	}

	_, _ = ctx.JSON(&StdResponse{
		Code:    CodeOK,
		Message: "删除成功~",
	})
}
