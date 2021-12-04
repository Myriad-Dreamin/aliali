package ali_drive

import (
	"errors"
	"github.com/go-resty/resty/v2"
	"io"
	"io/ioutil"
	"net/http"
)

var (
	ErrBodyIsNotStream    = errors.New("only workaround for handling big body stream")
	ErrBodyIsNotCloneable = errors.New("cloneable only workaround for handling big body stream")
)

type Cloneable interface {
	Clone() interface{}
}

type CloneableReader interface {
	io.Reader
	Cloneable
}

type CloneableReadCloser interface {
	io.ReadCloser
	Cloneable
}

func getCopyableReader(r *resty.Request) (Cloneable, error) {
	if cReader, ok := r.Body.(CloneableReader); ok {
		return cReader, nil
	}
	return nil, ErrBodyIsNotCloneable
}

func (y *Ali) workRoundUploadFile(c *resty.Client, r *resty.Request) (err error) {
	if r.Body == nil {
		return ErrBodyIsNotStream
	}
	if reader, ok := r.Body.(io.Reader); ok {
		r.RawRequest, err = http.NewRequest(r.Method, r.URL, reader)
	} else {
		return ErrBodyIsNotStream
	}

	if err != nil {
		return
	}

	// Assign close connection option
	// r.RawRequest.Close = c.closeConnection

	// Add headers into http request
	r.RawRequest.Header = r.Header

	// Add cookies from client instance into http request
	for _, cookie := range c.Cookies {
		r.RawRequest.AddCookie(cookie)
	}

	// Add cookies from request instance into http request
	for _, cookie := range r.Cookies {
		r.RawRequest.AddCookie(cookie)
	}

	// Use context if it was specified
	//if r.ctx != nil {
	//	req = r.RawRequest.WithContext(r.ctx)
	//}

	sink, err := getCopyableReader(r)
	if err != nil && err != ErrBodyIsNotCloneable {
		return err
	}

	// assign get body func for the underlying raw request instance
	r.RawRequest.GetBody = func() (io.ReadCloser, error) {
		if sink != nil {
			cReader := sink.Clone()
			if reader, ok := cReader.(io.Reader); ok {
				return ioutil.NopCloser(reader), nil
			}
			if reader, ok := cReader.(io.ReadCloser); ok {
				return reader, nil
			}
			return nil, ErrBodyIsNotStream
		}
		return nil, nil
	}

	return nil
}

func doRestyRequest(client *http.Client, req *resty.Request) (*resty.Response, error) {
	resp, err := client.Do(req.RawRequest)

	response := &resty.Response{
		Request:     req,
		RawResponse: resp,
	}

	if err != nil {
		//response.setReceivedAt()
		return response, err
	}

	return response, nil
}

func (y *Ali) FileUploadPart(reqBody *ApiFileUploadPartRequest) *ApiFileUploadPartResponse {
	url := reqBody.Uri
	req := y.r(y.uploadClient).
		SetBody(reqBody.Reader)

	req.Method = "PUT"
	req.URL = url

	if err := y.workRoundUploadFile(y.uploadClient, req); err != nil {
		y.suppress.Suppress(err)
		return nil
	}

	var resp = new(ApiFileUploadPartResponse)
	if _, err := doRestyRequest(y.uploadClient.GetClient(), req); err != nil {
		y.suppress.Suppress(err)
		return nil
	}
	return resp
}
