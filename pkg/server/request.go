package server

import (
	"github.com/gorilla/mux"
	"net/http"
)

type request struct {
	Request        *http.Request
	PathParameters map[string]string
}

type Request interface {
	GetPathParameters() map[string]string
	GetRequest() *http.Request
}

func Wrap(r *http.Request) *request {
	return &request{
		Request:        r,
		PathParameters: mux.Vars(r),
	}
}

func (r *request) GetRequest() *http.Request {
	return r.Request
}

func (r *request) GetPathParameters() map[string]string {
	return r.PathParameters
}
