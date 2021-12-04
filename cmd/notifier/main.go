package main

import (
	"github.com/Myriad-Dreamin/aliali/dispatcher"
	ali_notifier "github.com/Myriad-Dreamin/aliali/pkg/ali-notifier"
	"github.com/Myriad-Dreamin/aliali/pkg/suppress"
	"github.com/Myriad-Dreamin/aliali/server"
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	var s = suppress.PanicAll{}

	notifier := &ali_notifier.HttpRecorderNotifier{
		CapturePath: "/rec",
	}

	var d = dispatcher.NewDispatcher(dispatcher.WithNotifier(notifier))

	notifier.StorePath = d.GetConfig().AliDrive.RootPath

	r := iris.New().Configure(iris.WithoutBanner)

	r.AllowMethods(iris.MethodOptions)
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"OPTIONS", "HEAD", "GET", "POST", "DELETE"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
	r.Use(crs)
	notifier.ExposeHttp(r)
	(&server.Server{
		S:  s,
		DB: d.GetDatabase(),
	}).ExposeHttp(r)
	go func() {
		s.Suppress(r.Listen(":10305"))
	}()
	go func() {
		_ = http.ListenAndServe("0.0.0.0:10306", nil)
	}()

	d.Loop()
}
