package registry

import (
	"bytes"
	context2 "context"
	"encoding/json"
	"github.com/Myriad-Dreamin/aliali/server"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"io/ioutil"
)

type PostIdentRequest struct {
	Server string `json:"server"`
}

func (reg *Registry) DelegatePostRequest(ctx *context.Context) {
	var req PostIdentRequest
	var reqBody *bytes.Buffer
	if b, err := ctx.GetBody(); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		reg.Logger.Println(err.Error())
		_, _ = ctx.JSON(server.StdResponse{
			Code:    server.CodeInvalidParams,
			Message: "解析请求参数发生错误，查看后台日志了解内容...",
		})
		return
	} else {
		err = json.Unmarshal(b, &req)
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			reg.Logger.Println(err.Error())
			_, _ = ctx.JSON(server.StdResponse{
				Code:    server.CodeInvalidParams,
				Message: "解析请求参数发生错误，查看后台日志了解内容...",
			})
			return
		}
		reqBody = bytes.NewBuffer(b)
	}

	reqRaw := ctx.Request().Clone(context2.TODO())
	if !reg.ModifyHost(ctx, req.Server, reqRaw) {
		return
	}
	reqRaw.Body = ioutil.NopCloser(reqBody)
	reg.Delegate(ctx, reqRaw)
}
