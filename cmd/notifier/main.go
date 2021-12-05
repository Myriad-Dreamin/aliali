package main

import (
	"github.com/Myriad-Dreamin/aliali/dispatcher"
	ali_notifier "github.com/Myriad-Dreamin/aliali/pkg/ali-notifier"
	iris_cors "github.com/Myriad-Dreamin/aliali/pkg/iris-cors"
	"github.com/Myriad-Dreamin/aliali/pkg/suppress"
	"github.com/Myriad-Dreamin/aliali/server"
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
	iris_cors.Use(r)
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
