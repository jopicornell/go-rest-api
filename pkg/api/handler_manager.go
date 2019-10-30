package api

import (
	"encoding/json"
	"github.com/jopicornell/go-rest-api/internals/errors"
	"net/http"
)

func HandlerManager(handler func(http.ResponseWriter, *http.Request) (interface{}, error)) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		resource, err := handler(w, r)
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
	case errors.InternalServerError:
		w.WriteHeader(500)
	case errors.NotFound:
		w.WriteHeader(404)
	}
}
