package registry

import (
	"fmt"
	ali_notifier "github.com/Myriad-Dreamin/aliali/pkg/ali-notifier"
	"github.com/Myriad-Dreamin/aliali/server"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/core/router"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Info struct {
	Server         string    `json:"server"`
	Name           string    `json:"name"`
	Host           string    `json:"host"`
	Port           string    `json:"port"`
	UploaderStatus uint64    `json:"uploaderStatus"`
	RecorderStatus uint64    `json:"recorderStatus"`
	LastActive     time.Time `json:"last_active"`
}

type GetRegistryResponse struct {
	Items []*Info `json:"items"`
}

type Registry struct {
	*server.Server
	liveList map[string]*Info
	Config   *ali_notifier.RegistryConfig
	Servers  map[string]*ali_notifier.RegistryConfig
	Ident    string
	m        sync.Mutex
}

func NewRegistry(srv *server.Server, cfg *ali_notifier.Config, ident string) *Registry {
	return &Registry{
		Server:   srv,
		liveList: make(map[string]*Info),
		Ident:    ident,
		Config:   cfg.Servers[ident],
		Servers:  cfg.Servers,
	}
}

type PostRegisterRequest struct {
	Name           string `json:"name"`
	Server         string `json:"server"`
	Port           string `json:"port"`
	Secret         string `json:"secret"`
	RecorderStatus uint64 `json:"recorderStatus"`
}

func (reg *Registry) PostRegister(ctx *context.Context) {
	var req PostRegisterRequest
	if !reg.ReadJSON(ctx, &req) {
		return
	}

	if reg.Config.Secret != req.Secret {
		return
	}

	var h = strings.Split(ctx.Host(), ":")[0]

	reg.m.Lock()
	if _, ok := reg.liveList[req.Server]; !ok {
		reg.Server.Logger.Printf("new server %s => %s:%s", req.Server, h, req.Port)
	}
	reg.liveList[req.Server] = &Info{
		Server:         req.Server,
		Host:           h,
		Port:           req.Port,
		Name:           req.Name,
		UploaderStatus: 1,
		RecorderStatus: req.RecorderStatus,
		LastActive:     time.Now(),
	}
	reg.m.Unlock()
	_, _ = ctx.JSON(server.StdResponse{
		Message: "收到啦~",
	})
}

func (reg *Registry) GetRegistry(ctx *context.Context) {
	var resp GetRegistryResponse
	reg.m.Lock()
	for _, v := range reg.liveList {
		resp.Items = append(resp.Items, v)
	}
	for k, v := range reg.Servers {
		if v.Upstream != reg.Ident {
			continue
		}
		if _, ok := reg.liveList[k]; ok {
			continue
		} else {
			resp.Items = append(resp.Items, &Info{
				Server:         k,
				Name:           v.Name,
				Host:           v.ServerHost,
				Port:           v.RegistryPort,
				RecorderStatus: 0,
				UploaderStatus: 0,
			})
		}
	}
	reg.m.Unlock()

	ctx.StatusCode(200)
	_, _ = ctx.JSON(resp)
}

func (reg *Registry) ExposeHttp(r *iris.Application) {
	r.Handle("GET", "ping", func(ctx *context.Context) {
		_, _ = ctx.JSON(&server.PingResponse{
			Version: "v1.0.0",
			Message: "notifier registry",
		})
	})
	r.Handle("POST", "registry/v1/register", reg.PostRegister)
	r.Handle("POST", "api/v1/login", reg.Server.Login)
	r.PartyFunc("api/v1", func(p router.Party) {
		p.Use(reg.JwtHandler().Serve)
		p.Handle("GET", "/uploads", reg.DelegateGetRequest)
		p.Handle("DELETE", "/upload", reg.DelegatePostRequest)
		p.Handle("GET", "/registry", reg.GetRegistry)
	})
}

func (reg *Registry) ModifyHost(ctx *context.Context, ident string, req *http.Request) bool {
	if info, ok := reg.liveList[ident]; ok {
		req.URL.Scheme = "http"
		var hp = fmt.Sprintf("%s:%s", info.Host, info.Port)
		req.URL.Host = hp
		req.RequestURI = req.URL.String()
		req.Host = hp
		req.RemoteAddr = ""
		return true
	} else {
		ctx.StatusCode(iris.StatusBadRequest)
		_, _ = ctx.JSON(server.StdResponse{
			Code:    server.CodeNoSuchServer,
			Message: fmt.Sprintf("请求的服务器未注册: %s...", ident),
		})
		return false
	}
}

func (reg *Registry) Loop() {
	dur := time.Second * 33
	tick := time.NewTicker(dur)
	for {
		select {
		case t := <-tick.C:
			var clearList []string
			reg.m.Lock()
			for k, v := range reg.liveList {
				if t.Sub(v.LastActive) > dur {
					reg.Server.Logger.Printf("remove server %s => %s:%s", v.Server, v.Host, v.Port)
					clearList = append(clearList, k)
				}
			}
			for x := range clearList {
				delete(reg.liveList, clearList[x])
			}
			reg.m.Unlock()
		}
	}
}
