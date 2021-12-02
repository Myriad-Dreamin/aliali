package ali_drive

//go:generate go run ./generate api.def.yaml

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Myriad-Dreamin/aliali/pkg/suppress"
	"github.com/go-resty/resty/v2"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	//  "yp/config"
)

type Ali struct {
	client       *resty.Client
	uploadClient *resty.Client
	suppress     suppress.ISuppress
	accessToken  string

	RefreshInfo ApiRefreshResponse `json:"refresh_info"`
	Headers     [][2]string        `json:"headers"`
}

func NewAli() *Ali {
	return &Ali{
		client: resty.New(),
		uploadClient: resty.New().SetPreRequestHook(func(client *resty.Client, request *http.Request) error {
			request.Header.Set("Content-Type", "")
			return nil
		}),
	}
}

func (y *Ali) SetAccessToken(s string) {
	y.accessToken = s
}

type DataItem struct {
	Category             string                 `json:"category"`
	ContentHash          string                 `json:"content_hash"`
	ContentHashName      string                 `json:"content_hash_name"`
	ContentType          string                 `json:"content_type"`
	Crc64Hash            string                 `json:"crc64_hash"`
	CreatedAt            string                 `json:"created_at"`
	DomainId             string                 `json:"domain_id"`
	DownloadUrl          string                 `json:"download_url"`
	DeviceId             string                 `json:"device_id"`
	EncryptMode          string                 `json:"encrypt_mode"`
	FileExtension        string                 `json:"file_extension"`
	FileId               string                 `json:"file_id"`
	Hidden               bool                   `json:"hidden"`
	Name                 string                 `json:"name"`
	ParentFileId         string                 `json:"parent_file_id"`
	PunishFlag           int                    `json:"punish_flag"`
	Size                 int                    `json:"size"`
	Starred              bool                   `json:"starred"`
	Status               string                 `json:"status"`
	Type                 string                 `json:"type"`
	UpdatedAt            string                 `json:"updated_at"`
	UploadId             string                 `json:"upload_id"`
	Url                  string                 `json:"url"`
	Thumbnail            string                 `json:"thumbnail"`
	VideoPreviewMetadata map[string]interface{} `json:"video_preview_metadata"`
}

var Yunpan = new(Ali)

func init() {
	//Yunpan.Refresh()
	Yunpan.Heartbeat()
}

func (y *Ali) r(c *resty.Client) *resty.Request {
	req := c.R()

	for i := range y.Headers {
		req.SetHeader(y.Headers[i][0], y.Headers[i][1])
	}

	return req
}

func Unmarshal(s suppress.ISuppress, d []byte, i interface{}) bool {
	err := json.Unmarshal(d, i)
	if err != nil {
		s.Suppress(err)
		return false
	}
	return true
}

func (y *Ali) processResp(res *resty.Response, err error) []byte {
	if err != nil {
		y.suppress.Suppress(err)
		return nil
	}
	if res.StatusCode() >= 300 || res.StatusCode() < 200 {
		y.suppress.Suppress(errors.New(string(res.Body())))
		return nil
	}

	return res.Body()
}

func (y *Ali) unmarshal(b []byte, i interface{}) bool {
	if b == nil {
		return false
	}

	err := json.Unmarshal(b, i)
	if err != nil {
		y.suppress.Suppress(err)
		return false
	}

	return true
}

func (y *Ali) setAuthHeader(req *resty.Request) {
	req.SetHeader("authorization", y.accessToken)
}

func (y *Ali) GetDownloadUrl(file_id string) (DataItem, error) {
	url := "https://api.aliyundrive.com/v2/file/get"
	data := map[string]interface{}{
		"drive_id":                y.RefreshInfo.DefaultDriveId,
		"file_id":                 file_id,
		"image_thumbnail_process": "image/resize,w_400/format,jpeg",
		"fields":                  "*",
		"image_url_process":       "image/resize,w_1920/format,jpeg",
		"order_by":                "updated_at",
		"order_direction":         "DESC",
		"video_thumbnail_process": "video/snapshot,t_0,f_jpg,ar_auto,w_300",
	}
	data_json, _ := json.Marshal(data)
	header := map[string]string{
		"Content-Type":  "application/json",
		"origin":        "https://www.aliyundrive.com",
		"referer":       "https://www.aliyundrive.com",
		"authorization": y.RefreshInfo.TokenType + " " + y.RefreshInfo.AccessToken,
	}
	respByte, _ := y.curl(url, "POST", string(data_json), header)
	res := DataItem{}
	err := json.Unmarshal([]byte(respByte), &res)
	return res, err
}

func (y *Ali) GetAudioPlayInfo(file_id string) (map[string]interface{}, error) {
	url := "https://api.aliyundrive.com/v2/databox/get_audio_play_info"
	data := map[string]interface{}{
		"drive_id": y.RefreshInfo.DefaultDriveId,
		"file_id":  file_id,
	}
	data_json, _ := json.Marshal(data)
	header := map[string]string{
		"Content-Type":  "application/json",
		"origin":        "https://www.aliyundrive.com",
		"referer":       "https://www.aliyundrive.com",
		"authorization": y.RefreshInfo.TokenType + " " + y.RefreshInfo.AccessToken,
	}
	respByte, _ := y.curl(url, "POST", string(data_json), header)
	info := map[string]interface{}{}
	err := json.Unmarshal([]byte(respByte), &info)
	return info, err
}

func (y *Ali) GetVideoPlayInfo(file_id string) (map[string]interface{}, error) {
	url := "https://api.aliyundrive.com/v2/file/get_video_preview_play_info"
	data := map[string]interface{}{
		"drive_id": y.RefreshInfo.DefaultDriveId,
		"file_id":  file_id,
		"category": "live_transcoding",
	}
	data_json, _ := json.Marshal(data)
	header := map[string]string{
		"Content-Type":  "application/json",
		"origin":        "https://www.aliyundrive.com",
		"referer":       "https://www.aliyundrive.com",
		"authorization": y.RefreshInfo.TokenType + " " + y.RefreshInfo.AccessToken,
	}
	respByte, _ := y.curl(url, "POST", string(data_json), header)
	info := map[string]interface{}{}
	err := json.Unmarshal([]byte(respByte), &info)
	return info, err
}

func (y *Ali) MultiDownloadUrl(data map[string]interface{}) (map[string]interface{}, error) {
	url := "https://api.aliyundrive.com/adrive/v1/file/multiDownloadUrl"
	data_json, _ := json.Marshal(data)
	header := map[string]string{
		"Content-Type":  "application/json",
		"origin":        "https://www.aliyundrive.com",
		"referer":       "https://www.aliyundrive.com",
		"authorization": y.RefreshInfo.TokenType + " " + y.RefreshInfo.AccessToken,
	}
	respByte, _ := y.curl(url, "POST", string(data_json), header)
	info := map[string]interface{}{}
	err := json.Unmarshal([]byte(respByte), &info)
	return info, err
}

func (y *Ali) curl(url string, options ...interface{}) ([]byte, error) {
	//options -》 method string,data string,hearder map[string]string
	//获取访问方法
	method := "GET"
	if options[0] != nil {
		method = options[0].(string)
	}
	//获取参数
	data := ""
	if options[1] != nil {
		data = options[1].(string)
	}
	//获取头
	header := map[string]string{}
	if options[2] != nil {
		header = options[2].(map[string]string)
	}
	req, _ := http.NewRequest(method, url, strings.NewReader(data))
	//设置请求头
	for key, value := range header {
		req.Header.Set(key, value)
	}
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (y *Ali) Heartbeat() {
	//var ch chan int
	//定时任务
	ticker := time.NewTicker(time.Second * 6500)
	go func() {
		for range ticker.C {
			fmt.Println("心跳启动")
			//执行
			//y.Refresh()
		}

	}()

}
