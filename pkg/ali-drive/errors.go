package ali_drive

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"runtime"
)

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
