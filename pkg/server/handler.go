package server

import (
	"net/http"
)

type HandlerFunc func(Request)

type Handler struct {
	http.Handler
	routes map[string]func()
}

func HandleHTTP(handler HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(NewRequest(r, w))
	}
}
