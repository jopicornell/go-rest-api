package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jopicornell/go-rest-api/internals/models"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

var UserContextKey = "user"

type request struct {
	Request        *http.Request
	PathParameters map[string]string
	user           *models.User
}

type Request interface {
	GetPathParameters() map[string]string
	GetParamUInt(param string) uint
	GetRequest() *http.Request
	GetUser() *models.User
	GetBody() []byte
	GetBodyMarshalled(interface{})
}

func NewRequest(req *http.Request) *request {
	return &request{
		Request:        req,
		PathParameters: mux.Vars(req),
	}
}

func (r *request) GetRequest() *http.Request {
	return r.Request
}

func (r *request) GetUser() *models.User {
	user := r.Request.Context().Value(UserContextKey).(models.User)
	return &user
}

func (r *request) GetBody() []byte {
	if body, err := ioutil.ReadAll(r.Request.Body); err == nil {
		return body
	} else {
		panic(Error{StatusCode: http.StatusBadRequest})
	}
}

func (r *request) GetBodyMarshalled(ifc interface{}) {
	if err := json.Unmarshal(r.GetBody(), ifc); err != nil {
		log.Println(fmt.Sprintf("error unmarshalling: %s", err))
		panic(Error{StatusCode: http.StatusBadRequest})
	}
}

func (r *request) GetPathParameters() map[string]string {
	return r.PathParameters
}

func (r *request) GetParam(param string) string {
	return r.PathParameters[param]
}

func (r *request) GetParamUInt(param string) uint {
	if value, err := strconv.Atoi(r.PathParameters[param]); err == nil {
		return uint(value)
	} else {
		panic(err)
	}
}
