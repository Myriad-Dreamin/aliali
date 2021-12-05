package ali_drive

//go:generate go run ./generate api.def.yaml

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Myriad-Dreamin/aliali/pkg/suppress"
	"github.com/go-resty/resty/v2"
	"net/http"
	"runtime"
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

type AliHttpError struct {
	// you should know that it is *resty.Response if you want to access it
	Response   interface{}
	Caller     string
	StatusCode int
	Message    string
}

func (a *AliHttpError) Error() string {
	return fmt.Sprintf("%s: Http Error %d, Message %s", a.Caller, a.StatusCode, a.Message)
}

type AliSemaError struct {
	// you should know that it is *resty.Response if you want to access it
	Response   interface{}
	Caller     string
	StatusCode int
	*ApiErrorResponse
}

func (a *AliSemaError) Error() string {
	return fmt.Sprintf("%s: %s(%d): %s", a.Caller, a.Code, a.StatusCode, a.ApiErrorResponse.Message)
}

func httpError(caller string, res *resty.Response) error {
	return &AliHttpError{
		Response:   res,
		Caller:     caller,
		StatusCode: res.StatusCode(),
		Message:    res.Status(),
	}
}

func semaError(caller string, res *resty.Response, msg *ApiErrorResponse) error {
	return &AliSemaError{
		Response:         res,
		Caller:           caller,
		StatusCode:       res.StatusCode(),
		ApiErrorResponse: msg,
	}
}

func getCaller(skip int) string {
	pc, _, _, ok := runtime.Caller(skip + 1)
	if ok {
		details := runtime.FuncForPC(pc)
		if details != nil {
			return details.Name()
		}
	}
	return ""
}

func (y *Ali) InjectSemaError(statusCode int, msg *ApiErrorResponse) {
	y.suppress.Suppress(&AliSemaError{
		Response:         nil,
		Caller:           "",
		StatusCode:       statusCode,
		ApiErrorResponse: msg,
	})
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
		if sErr != nil || len(messageUnpack.Code) == 0 {
			y.suppress.Suppress(semaError(getCaller(1), res, &messageUnpack))
		} else {
			y.suppress.Suppress(errors.New(string(res.Body())))
		}
		return nil
	}

	return res.Body()
}
