package registry

import (
	"github.com/Myriad-Dreamin/aliali/server"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"io"
	"net/http"
	"strings"
)

func (reg *Registry) Delegate(ctx *context.Context, req *http.Request) {
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		reg.Logger.Println(err.Error())
		_, _ = ctx.JSON(server.StdResponse{
			Code:    server.CodeErr,
			Message: "重定向请求发生错误，查看后台日志了解内容...",
		})
		return
	}

	ctx.StatusCode(resp.StatusCode)
	for i := range resp.Header {
		ctx.Header(i, "")
		ctx.Header(i, strings.Join(resp.Header[i], " "))
	}

	_ = ctx.StreamWriter(func(w io.Writer) error {
		_, err := io.Copy(w, resp.Body)
		if err != nil {
			reg.Logger.Println(err.Error())
		}
		_ = resp.Body.Close()
		return io.EOF
	})
	ctx.StopExecution()
	_ = resp.Body.Close()
}
