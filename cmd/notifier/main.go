package main

import (
	"github.com/Myriad-Dreamin/aliali/dispatcher"
	ali_notifier "github.com/Myriad-Dreamin/aliali/pkg/ali-notifier"
	"github.com/Myriad-Dreamin/aliali/pkg/suppress"
	"github.com/kataras/iris/v12"
)

func main() {
	var s = suppress.PanicAll{}

	notifier := &ali_notifier.HttpRecorderNotifier{
		CapturePath: "/rec",
	}

	var d = dispatcher.NewDispatcher(dispatcher.WithNotifier(notifier))

	notifier.StorePath = d.GetConfig().AliDrive.RootPath

	r := iris.New().Configure(iris.WithoutBanner)

	notifier.ExposeHttp(r)
	go func() {
		s.Suppress(r.Listen(":10305"))
	}()

	d.Loop()
}
