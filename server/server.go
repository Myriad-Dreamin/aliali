package server

import (
	ali_notifier "github.com/Myriad-Dreamin/aliali/pkg/ali-notifier"
	"github.com/Myriad-Dreamin/aliali/pkg/suppress"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/core/router"
	"gorm.io/gorm"
	"log"
)

const (
	CodeOK int = iota
	CodeErr
	CodeInvalidParams
	CodeDBExecutionError
	CodeLoginNotConfigured
	CodeLoginNoSuchAccount
	CodeLoginWrongPassword
	CodeNoSuchId
	CodeNoSuchServer
)

type StdResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type PingResponse struct {
	Version string `json:"version"`
	Message string `json:"message"`
}

type Server struct {
	Logger *log.Logger
	S      suppress.ISuppress
	DB     *gorm.DB
	Config *ali_notifier.BackendConfig
}

func (srv *Server) ExposeHttp(r *iris.Application) {

	r.Handle("GET", "ping", func(ctx *context.Context) {
		_, _ = ctx.JSON(&PingResponse{
			Version: "v1.0.0",
			Message: "notifier backend",
		})
	})

	r.Handle("POST", "api/v1/login", srv.Login)

	r.PartyFunc("api/v1", func(p router.Party) {
		p.Use(srv.JwtHandler().Serve)
		p.Handle("GET", "/uploads", srv.GetUploadList)
		p.Handle("DELETE", "/upload", srv.DeleteUpload)
	})
}

func (srv *Server) ReadJSON(ctx *context.Context, req interface{}) bool {
	if err := ctx.ReadJSON(req); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		srv.Logger.Println(err.Error())
		_, _ = ctx.JSON(StdResponse{
			Code:    CodeInvalidParams,
			Message: "解析请求参数发生错误，查看后台日志了解内容...",
		})
		return false
	}
	return true
}

func (srv *Server) ReadQuery(ctx *context.Context, req interface{}) bool {
	if err := ctx.ReadQuery(req); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		srv.Logger.Println(err.Error())
		_, _ = ctx.JSON(StdResponse{
			Code:    CodeInvalidParams,
			Message: "解析请求参数发生错误，查看后台日志了解内容...",
		})
		return false
	}
	return true
}
