package servertesting

import (
	"encoding/json"
	"github.com/jopicornell/go-rest-api/pkg/server"
	"net/http"
	"net/http/httptest"
)

type ResponseMock struct {
	httptest.ResponseRecorder
	errorCalled *server.Error
}

// constructor to get a new response
func NewResponse() *ResponseMock {
	return &ResponseMock{ResponseRecorder: *httptest.NewRecorder()}
}

// get the writer so you can handle responses in any other way
func (r *ResponseMock) GetWriter() http.ResponseWriter {
	return r
}

//responds with the parsed interface to JSON
func (r *ResponseMock) RespondJSON(statusCode int, ifc interface{}) {
	r.ResponseRecorder.Header().Add("Content-Type", "application/json")
	r.ResponseRecorder.WriteHeader(statusCode)
	if err := json.NewEncoder(r).Encode(ifc); err == nil {
	} else {
		panic(err)
	}
}

// responds without a body
func (r *ResponseMock) Respond(statusCode int) {
	r.ResponseRecorder.WriteHeader(statusCode)
	if _, err := r.ResponseRecorder.Write([]byte{}); err != nil {
		panic(err)
	}
}

// sets the "ErrorHasBeenCalled" flag so you can know, without panicking, if an error has been thrown
func (r *ResponseMock) Error(error *server.Error) {
	if r.errorCalled != nil {
		panic("method Error has been called twice")
	}
	r.errorCalled = error
	r.Code = error.StatusCode
}

// return the error called or nil if it was not called
func (r *ResponseMock) ErrorCalled() *server.Error {
	return r.errorCalled
}
