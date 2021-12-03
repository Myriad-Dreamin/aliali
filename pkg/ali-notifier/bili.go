package ali_notifier

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

type BiliRecorderNotifier struct {
	Notifier
}

type BiliNotificationRequest struct {
}

type BiliNotificationResponse struct {
	Message string `json:"message"`
}

func (b *BiliRecorderNotifier) Run() error {
	r := iris.New() // .Configure(iris.WithoutBanner)
	r.Handle("POST", "/notifier/bilibili", func(ctx context.Context) {
		var req BiliNotificationRequest
		if err := ctx.ReadJSON(&req); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			_, _ = ctx.WriteString(err.Error())
			return
		}

		_, _ = ctx.JSON(BiliNotificationResponse{
			Message: "收到啦~",
		})
	})
	return r.Listen(":10305")
}
