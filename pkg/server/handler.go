package server

import (
	"net/http"
)

type HandlerFunc func(Response, Request)

type handler struct{}

type Handler interface {
}

func HandleHTTP(handler HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(NewResponse(w), NewRequest(r))
	}
}
