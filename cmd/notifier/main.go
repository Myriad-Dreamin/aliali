package main

import (
	"fmt"
	"github.com/Myriad-Dreamin/aliali/dispatcher"
	ali_notifier "github.com/Myriad-Dreamin/aliali/pkg/ali-notifier"
	iris_cors "github.com/Myriad-Dreamin/aliali/pkg/iris-cors"
	"github.com/Myriad-Dreamin/aliali/pkg/suppress"
	"github.com/Myriad-Dreamin/aliali/registry"
	"github.com/Myriad-Dreamin/aliali/server"
	"github.com/kataras/iris/v12"
	"gorm.io/gorm"
	"log"
	_ "net/http/pprof"
	"os"
	"strings"
)

func getRegistryHostPort(ident string, cfg *ali_notifier.Config) (string, string) {
	var host = "0.0.0.0"
	var port = "10308"
	if len(ident) != 0 && cfg.Servers != nil && cfg.Servers[ident] != nil {
		if len(cfg.Servers[ident].RegistryHost) != 0 {
			host = cfg.Servers[ident].RegistryHost
		}
		if len(cfg.Servers[ident].RegistryPort) != 0 {
			port = cfg.Servers[ident].RegistryPort
		}
	}

	return host, port
}

func createClient(ident string, cfg *ali_notifier.Config) *registry.Client {
	if len(ident) != 0 && cfg.Servers != nil && cfg.Servers[ident] != nil && len(cfg.Servers[ident].Upstream) != 0 {
		if cfg.Servers[cfg.Servers[ident].Upstream] == nil {
			return nil
		}
		cfg.Servers[ident].UpstreamSecret = cfg.Servers[cfg.Servers[ident].Upstream].Secret
		cfg.Servers[ident].UpstreamHost, cfg.Servers[ident].UpstreamPort =
			getRegistryHostPort(cfg.Servers[ident].Upstream, cfg)
		return registry.NewClient(ident, cfg.Servers[ident])
	}
	return nil
}

func clientMain(c *registry.Client) {
	_ = c.Run()
}

func createServer(s suppress.ISuppress, cfg *ali_notifier.BackendConfig, db *gorm.DB) *server.Server {
	return &server.Server{
		Logger: log.New(os.Stderr, "[backend] ", log.Llongfile|log.LUTC),
		S:      s,
		DB:     db,
		Config: cfg,
	}
}

func backendListen(ident string, cfg *ali_notifier.Config, r *iris.Application) error {
	var host = "0.0.0.0"
	var port = "10305"
	if len(ident) != 0 && cfg.Servers != nil && cfg.Servers[ident] != nil {
		if len(cfg.Servers[ident].ServerHost) != 0 {
			host = cfg.Servers[ident].ServerHost
		}
		if len(cfg.Servers[ident].ServerPort) != 0 {
			port = cfg.Servers[ident].ServerPort
		}
	}

	return r.Listen(fmt.Sprintf("%s:%s", host, port))
}

func registryListen(ident string, cfg *ali_notifier.Config, r *iris.Application) error {
	var host, port = getRegistryHostPort(ident, cfg)
	return r.Listen(fmt.Sprintf("%s:%s", host, port))
}

func notifierMain(ident string) {
	var s = suppress.PanicAll{}

	notifier := &ali_notifier.HttpRecorderNotifier{
		CapturePath: "/rec",
	}

	var d = dispatcher.NewDispatcher(dispatcher.WithNotifier(notifier))

	notifier.StorePath = d.GetConfig().AliDrive.RootPath
	var c = createClient(ident, d.GetConfig())

	r := iris.New().Configure(iris.WithoutBanner)
	iris_cors.Use(r)
	notifier.ExposeHttp(r)
	createServer(s, &d.GetConfig().Backend, d.GetDatabase()).ExposeHttp(r)
	go func() {
		s.Suppress(backendListen(ident, d.GetConfig(), r))
	}()
	if len(ident) != 0 {
		go clientMain(c)
	}
	d.Loop()
}

func backendMain(ident string) {
	var s = suppress.PanicAll{}
	var dm = &dispatcher.DBManager{S: s}
	var cm = &dispatcher.ConfigManager{S: s}
	var db = dm.OpenSqliteDB("deployment/backend-workdir/ali.db")
	var cfg = cm.ReadConfig(dispatcher.DefaultConfigPath)
	var srv = createServer(s, &cfg.Backend, db)
	var c = createClient(ident, cfg)

	r := iris.New()
	iris_cors.Use(r)
	srv.ExposeHttp(r)

	if c != nil {
		go clientMain(c)
	}
	s.Suppress(backendListen(ident, cfg, r))
}

func registryMain(ident string) {
	var s = suppress.PanicAll{}
	var cm = &dispatcher.ConfigManager{S: s}
	var cfg = cm.ReadConfig(dispatcher.DefaultConfigPath)
	var srv = registry.NewRegistry(createServer(s, &cfg.Backend, nil), cfg, ident)
	var c = createClient(ident, cfg)

	r := iris.New()
	iris_cors.Use(r)
	srv.ExposeHttp(r)
	if c != nil {
		go clientMain(c)
	}
	go func() {
		s.Suppress(registryListen(ident, cfg, r))
	}()

	srv.Loop()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("starts with notifier standalone mode")
		notifierMain("")
	} else {
		mode := os.Args[1]
		ms := strings.SplitN(mode, ":", 2)
		mode = ms[0]
		ident := ms[0]
		if len(ms) > 1 {
			ident = ms[1]
		}

		fmt.Printf("starts with %s mode\n", mode)
		switch mode {
		case "registry":
			registryMain(ident)
		case "backend":
			backendMain(ident)
		case "notifier":
			notifierMain(ident)
		}
	}
}
