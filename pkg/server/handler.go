package server

import (
	"net/http"
)

type HandlerFunc func(Context)

type Handler struct {
	http.Handler
	routes map[string]func()
}

func HandleHTTP(handler HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(NewContext(r, w))
	}
}
