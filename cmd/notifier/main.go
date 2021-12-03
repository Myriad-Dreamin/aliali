package main

import (
	"github.com/Myriad-Dreamin/aliali/dispatcher"
	ali_notifier "github.com/Myriad-Dreamin/aliali/pkg/ali-notifier"
	"github.com/Myriad-Dreamin/aliali/pkg/suppress"
)

func main() {
	var s = suppress.PanicAll{}

	notifier := &ali_notifier.HttpRecorderNotifier{
		CapturePath: "/rec",
	}

	var d = dispatcher.NewDispatcher(dispatcher.WithNotifier(notifier))

	notifier.StorePath = d.GetConfig().AliDrive.RootPath

	go func() {
		s.Suppress(notifier.Run())
	}()

	d.Loop()
}
