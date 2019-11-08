package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
)

type context struct {
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	PathParameters map[string]string
}

type Context interface {
	GetPathParameters() map[string]string
	GetParamUInt(param string) uint
	GetRequest() *http.Request
	GetWriter() http.ResponseWriter
	GetBody() []byte
	GetBodyMarshalled(interface{})
	RespondJSON(statusCode int, ifc interface{})
	Respond(statusCode int)
}

func NewContext(req *http.Request, resWriter http.ResponseWriter) *context {
	return &context{
		Request:        req,
		ResponseWriter: resWriter,
		PathParameters: mux.Vars(req),
	}
}

func (r *context) GetRequest() *http.Request {
	return r.Request
}

func (r *context) GetWriter() http.ResponseWriter {
	return r.ResponseWriter
}

func (r *context) GetBody() []byte {
	if body, err := ioutil.ReadAll(r.Request.Body); err == nil {
		return body
	} else {
		r.Respond(http.StatusBadRequest)
		return nil
	}
}

func (r *context) GetBodyMarshalled(ifc interface{}) {
	if err := json.Unmarshal(r.GetBody(), ifc); err != nil {
		r.Respond(http.StatusBadRequest)
	}
}

func (r *context) GetPathParameters() map[string]string {
	return r.PathParameters
}

func (r *context) GetParam(param string) string {
	return r.PathParameters[param]
}

func (r *context) GetParamUInt(param string) uint {
	if value, err := strconv.Atoi(r.PathParameters[param]); err == nil {
		return uint(value)
	} else {
		panic(err)
	}
}

func (r *context) RespondJSON(statusCode int, ifc interface{}) {
	r.ResponseWriter.Header().Add("Content-Type", "application/json")
	r.ResponseWriter.WriteHeader(statusCode)
	if err := json.NewEncoder(r.ResponseWriter).Encode(ifc); err == nil {
	} else {
		panic(err)
	}
}

func (r *context) Respond(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	if _, err := r.ResponseWriter.Write([]byte{}); err != nil {
		panic(err)
	}
}
