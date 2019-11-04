package servertesting

import "net/http"

func New(r *http.Request, pathParameters map[string]string) *RequestMock {
	return &RequestMock{
		Request:        r,
		PathParameters: pathParameters,
	}
}

type RequestMock struct {
	Request        *http.Request
	PathParameters map[string]string
}

func (r *RequestMock) GetRequest() *http.Request {
	return r.Request
}

func (r *RequestMock) GetPathParameters() map[string]string {
	return r.PathParameters
}
