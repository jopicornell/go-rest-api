package servertesting

import (
	"io/ioutil"
	"net/http"
)

func NewRequest(r *http.Request, pathParameters map[string]string) *RequestMock {
	return &RequestMock{
		Request:        r,
		PathParameters: pathParameters,
	}
}

type RequestMock struct {
	Request        *http.Request
	PathParameters map[string]string
	ThrowError     error
}

func (r *RequestMock) GetRequest() *http.Request {
	return r.Request
}

func (r *RequestMock) GetBody() ([]byte, error) {
	if r.ThrowError != nil {
		return nil, r.ThrowError
	}
	body, err := ioutil.ReadAll(r.Request.Body)
	return body, err
}

func (r *RequestMock) GetPathParameters() map[string]string {
	return r.PathParameters
}
