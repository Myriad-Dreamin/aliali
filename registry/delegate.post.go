package registry

import (
	context2 "context"
	"github.com/kataras/iris/v12/context"
)

type PostIdentRequest struct {
	Server string `json:"server"`
}

func (reg *Registry) DelegatePostRequest(ctx *context.Context) {
	var req PostIdentRequest
	if !reg.ReadJSON(ctx, &req) {
		return
	}

	reqRaw := ctx.Request().Clone(context2.TODO())
	if !reg.ModifyHost(ctx, req.Server, reqRaw) {
		return
	}
	reg.Delegate(ctx, reqRaw)
}
