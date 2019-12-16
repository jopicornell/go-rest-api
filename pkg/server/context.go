package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jopicornell/go-rest-api/api/users/models"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

var UserContextKey = "user"

type innerContext struct {
	Request        *http.Request
	PathParameters map[string]string
	user           *models.User
	server         Server
}

type Context interface {
	GetPathParameters() map[string]string
	GetParamUInt(param string) uint
	GetParamInt(param string) int
	GetParamString(param string) string
	GetRequest() *http.Request
	GetUser() interface{}
	SetUser(user interface{})
	GetServer() Server
	GetBody() []byte
	GetBodyMarshalled(interface{})
}

func NewContext(req *http.Request, server Server) *innerContext {
	return &innerContext{
		Request:        req,
		PathParameters: mux.Vars(req),
		server:         server,
	}
}

func (r *innerContext) GetRequest() *http.Request {
	return r.Request
}

func (r *innerContext) GetServer() Server {
	return r.server
}

func (r *innerContext) GetUser() interface{} {
	user := r.Request.Context().Value(UserContextKey)
	return user
}

func (r *innerContext) SetUser(user interface{}) {
	newContext := context.WithValue(r.Request.Context(), UserContextKey, user)
	r.Request = r.Request.WithContext(newContext)
}

func (r *innerContext) SetKey(key string, value interface{}) {
	newContext := context.WithValue(r.Request.Context(), key, value)
	r.Request = r.Request.WithContext(newContext)
}

func (r *innerContext) GetKey(key string, value interface{}) {
	value = r.Request.Context().Value(key)
}

func (r *innerContext) GetBody() []byte {
	if body, err := ioutil.ReadAll(r.Request.Body); err == nil {
		return body
	} else {
		panic(&Error{StatusCode: http.StatusBadRequest})
	}
}

func (r *innerContext) GetBodyMarshalled(ifc interface{}) {
	if err := json.Unmarshal(r.GetBody(), ifc); err != nil {
		log.Println(fmt.Sprintf("error unmarshalling: %s", err))
		panic(&Error{StatusCode: http.StatusBadRequest})
	}
}

func (r *innerContext) GetPathParameters() map[string]string {
	return r.PathParameters
}

func (r *innerContext) GetParam(param string) string {
	return r.PathParameters[param]
}

func (r *innerContext) GetParamUInt(param string) uint {
	if value, err := strconv.Atoi(r.PathParameters[param]); err == nil {
		return uint(value)
	} else {
		panic(err)
	}
}

func (r *innerContext) GetParamString(param string) string {
	return r.PathParameters[param]
}

func (r *innerContext) GetParamInt(param string) int {
	if value, err := strconv.Atoi(r.PathParameters[param]); err == nil {
		return value
	} else {
		panic(err)
	}
}
