package ali_drive

//go:generate go run ./generate api.def.yaml

import (
	"encoding/json"
	"errors"
	"github.com/Myriad-Dreamin/aliali/pkg/suppress"
	"github.com/go-resty/resty/v2"
	"net/http"
)

type Ali struct {
	client       *resty.Client
	uploadClient *resty.Client
	suppress     suppress.ISuppress
	accessToken  string
	Headers      [][2]string `json:"headers"`
}

func NewAli(suppress suppress.ISuppress) *Ali {
	return &Ali{
		client:   resty.New(),
		suppress: suppress,
		uploadClient: resty.New().SetPreRequestHook(func(client *resty.Client, request *http.Request) error {
			request.Header.Set("Content-Type", "")
			return nil
		}),
	}
}

func (y *Ali) SetAccessToken(s string) {
	y.accessToken = s
}

func (y *Ali) GetAccessToken() string {
	return y.accessToken
}

func (y *Ali) r(c *resty.Client) *resty.Request {
	req := c.R()

	for i := range y.Headers {
		req.SetHeader(y.Headers[i][0], y.Headers[i][1])
	}

	return req
}

func (y *Ali) setAuthHeader(req *resty.Request) {
	req.SetHeader("authorization", y.accessToken)
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

func (y *Ali) processResp(res *resty.Response, err error) []byte {
	if err != nil {
		y.suppress.Suppress(err)
		return nil
	}
	if res.StatusCode() >= 300 || res.StatusCode() < 200 {
		var b = res.Body()
		if len(b) == 0 {
			y.suppress.Suppress(httpError(getCaller(1), res))
			return nil
		}
		var messageUnpack ApiErrorResponse
		sErr := json.Unmarshal(b, &messageUnpack)
		if sErr == nil && len(messageUnpack.Code) != 0 {
			y.suppress.Suppress(semaError(getCaller(1), res, &messageUnpack))
		} else {
			y.suppress.Suppress(errors.New(string(res.Body())))
		}
		return nil
	}

	return res.Body()
}
