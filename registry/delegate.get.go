package registry

import (
	context2 "context"
	"github.com/kataras/iris/v12/context"
)

type GetIdentRequest struct {
	Server string `url:"server"`
}

func (reg *Registry) DelegateGetRequest(ctx *context.Context) {
	var req GetIdentRequest
	if !reg.ReadQuery(ctx, &req) {
		return
	}

	reqRaw := ctx.Request().Clone(context2.TODO())
	if !reg.ModifyHost(ctx, req.Server, reqRaw) {
		return
	}
	reg.Delegate(ctx, reqRaw)
}
