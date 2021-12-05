package ali_notifier

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"log"
	"path/filepath"
	"strconv"
	"strings"
)

type HttpRecorderNotifier struct {
	Notifier
	CapturePath string
	StorePath   string
}

type BiliNotificationRequest struct {
	EventType      string                 `json:"EventType"`
	EventTimestamp string                 `json:"EventTimestamp"`
	EventId        string                 `json:"EventId"`
	EventData      map[string]interface{} `json:"EventData"`
}

type BiliNotificationResponse struct {
	Message string `json:"message"`
}

type WebhookNotificationRequest struct {
	Group      string `url:"group"`
	Path       string `url:"path"`
	RemotePath string `url:"remote"`
}

type WebhookNotificationResponse struct {
	Message string `json:"message"`
}

func SecureJoin(root string, components string) string {
	root = filepath.Clean(root)
	var res = filepath.Clean(filepath.Join(root, components))

	if !strings.HasPrefix(res, root) {
		return ""
	}

	return res
}

func (b *HttpRecorderNotifier) NotifyMaybeUploadEvent(ctx *context.Context, group string, relLocalPath string, remotePath string) {
	if len(relLocalPath) == 0 {
		ctx.StatusCode(iris.StatusBadRequest)
		_, _ = ctx.JSON(WebhookNotificationResponse{
			Message: "解析目标文件地址好像是空的...",
		})
		return
	}

	if len(remotePath) == 0 {
		ctx.StatusCode(iris.StatusBadRequest)
		_, _ = ctx.JSON(WebhookNotificationResponse{
			Message: "解析上传地址好像是空的...",
		})
		return
	}

	var localPath = SecureJoin(b.CapturePath, relLocalPath)

	if len(localPath) == 0 {
		ctx.StatusCode(iris.StatusBadRequest)
		_, _ = ctx.JSON(WebhookNotificationResponse{
			Message: "在访问其他文件夹吗？不要干坏事哦",
		})
		return
	}

	b.Emit(&FsUploadRequest{
		Group:      group,
		DriveID:    "",
		RemotePath: remotePath,
		LocalPath:  localPath,
	})
	_, _ = ctx.JSON(WebhookNotificationResponse{
		Message: "收到啦~",
	})
}

func (b *HttpRecorderNotifier) NotifyBilibiliEvent(ctx *context.Context) {
	var req BiliNotificationRequest
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		fmt.Println(err.Error())
		_, _ = ctx.JSON(BiliNotificationResponse{
			Message: "解析请求参数发生错误，查看后台日志了解内容...",
		})
		return
	}

	if len(req.EventType) != 0 {
		log.Printf("接收到事件 %s\n", req.EventType)

		if req.EventType == "FileClosed" {
			var relPath string
			var group string
			if req.EventData != nil {
				if r, ok := req.EventData["RelativePath"]; ok {
					if r2, ok := r.(string); ok {
						relPath = r2
					}
				}
				if r, ok := req.EventData["RoomId"]; ok {
					if r2, ok := r.(float64); ok {
						group = strconv.FormatInt(int64(r2), 10)
						if r3, ok := req.EventData["Name"]; ok {
							if r4, ok := r3.(string); ok {
								group = fmt.Sprintf("%s-%s", group, r4)
							}
						}
					}
				}
			}

			if len(relPath) == 0 {
				log.Printf("诶，文件关闭事件里没有文件路径信息... %s %v\n", req.EventTimestamp, req.EventData)
				ctx.StatusCode(iris.StatusBadRequest)
				_, _ = ctx.JSON(BiliNotificationResponse{
					Message: "文件关闭事件里没有文件路径信息...",
				})
				return
			}

			log.Printf("检测到待上传的文件 %s\n", relPath)
			b.NotifyMaybeUploadEvent(ctx, group, relPath, filepath.Join(b.StorePath, relPath))
			return
		}
	}

	_, _ = ctx.JSON(BiliNotificationResponse{
		Message: "收到啦~",
	})
}

func (b *HttpRecorderNotifier) NotifyWebhookEvent(ctx *context.Context) {
	var req WebhookNotificationRequest
	if err := ctx.ReadQuery(&req); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		fmt.Println(err.Error())
		_, _ = ctx.JSON(WebhookNotificationResponse{
			Message: "解析请求参数发生错误，查看后台日志了解内容...",
		})
		return
	}

	b.NotifyMaybeUploadEvent(ctx, req.Group, req.Path, req.RemotePath)
}

func (b *HttpRecorderNotifier) ExposeHttp(r *iris.Application) {
	r.Handle("POST", "/notifier/bilibili", b.NotifyBilibiliEvent)
	r.Handle("GET", "/notifier/webhook", b.NotifyWebhookEvent)
}
