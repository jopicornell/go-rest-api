package server

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type response struct {
	responseWriter http.ResponseWriter
}

type Response interface {
	GetWriter() http.ResponseWriter
	RespondJSON(statusCode int, ifc interface{})
	RespondValidationErrors(statusCode int, errors validator.ValidationErrors)
	Respond(statusCode int)
	Error(error *Error)
}

// constructor to get a new response
func NewResponse(w http.ResponseWriter) Response {
	return &response{responseWriter: w}
}

// get the writer so you can handle responses in any other way
func (r *response) GetWriter() http.ResponseWriter {
	return r.responseWriter
}

//responds with the parsed interface to JSON
func (r *response) RespondJSON(statusCode int, ifc interface{}) {
	r.responseWriter.Header().Add("Content-Type", "application/json")
	r.responseWriter.WriteHeader(statusCode)
	if err := json.NewEncoder(r.responseWriter).Encode(ifc); err != nil {
		panic(err)
	}
}

//responds with the parsed interface to JSON
func (r *response) RespondValidationErrors(statusCode int, errors validator.ValidationErrors) {
	r.responseWriter.Header().Add("Content-Type", "application/json")
	r.responseWriter.WriteHeader(statusCode)
	validationErrors := make(map[string]map[string]string)
	for _, validatorError := range errors {
		if validationErrors[validatorError.Namespace()] == nil {
			validationErrors[validatorError.Namespace()] = make(map[string]string)
		}
		validationErrors[validatorError.Namespace()][validatorError.Tag()] = validatorError.Translate(nil)
	}
	if err := json.NewEncoder(r.responseWriter).Encode(validationErrors); err != nil {
		panic(err)
	}
}

// responds without a body
func (r *response) Respond(statusCode int) {
	r.responseWriter.WriteHeader(statusCode)
	if _, err := r.responseWriter.Write([]byte{}); err != nil {
		panic(err)
	}
}

// panics an error that can be handled and recovered by the server
func (r *response) Error(error *Error) {
	panic(error)
}
