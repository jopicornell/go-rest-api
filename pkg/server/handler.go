package server

import (
	"net/http"
)

type HandlerFunc func(Response, Request)

type Handler interface {
	http.Handler
	ConfigureRoutes()
}

func HandleHTTP(handler HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(NewResponse(w), NewRequest(r))
	}
}
