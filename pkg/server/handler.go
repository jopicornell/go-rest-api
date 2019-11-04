package server

import (
	"encoding/json"
	"github.com/jopicornell/go-rest-api/internals/errors"
	"net/http"
)

type HandlerSerializer func(func(http.ResponseWriter, Request) (interface{}, error)) http.HandlerFunc

type HandlerFunc func(http.ResponseWriter, Request) (interface{}, error)

type Handler struct {
	*Server
	routes map[string]func()
}

func HandleJSONResponse(handler HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resource, err := handler(w, Wrap(r))
		if err != nil {
			handleError(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(resource); err != nil {

		}
	}
}

func handleError(w http.ResponseWriter, err error) {
	switch err {
	case errors.NotFound:
		w.WriteHeader(404)
	default:
		w.WriteHeader(500)
	}
}
