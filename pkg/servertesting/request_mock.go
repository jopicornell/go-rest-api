package servertesting

import (
	"encoding/json"
	"github.com/jopicornell/go-rest-api/internals/models"
	"github.com/jopicornell/go-rest-api/pkg/server"
	"gopkg.in/guregu/null.v3"
	"io/ioutil"
	"net/http"
	"strconv"
)

func NewRequest(r *http.Request, pathParameters map[string]string) *RequestMock {
	return &RequestMock{
		Request:        r,
		PathParameters: pathParameters,
	}
}

type RequestMock struct {
	Request        *http.Request
	PathParameters map[string]string
	ResponseWriter http.ResponseWriter
	ThrowError     error
	User           *models.User
}

func (r *RequestMock) GetRequest() *http.Request {
	return r.Request
}

func (r *RequestMock) GetWriter() http.ResponseWriter {
	return r.ResponseWriter
}

func (r *RequestMock) GetBody() []byte {
	if body, err := ioutil.ReadAll(r.Request.Body); err == nil {
		return body
	} else {
		panic(server.Error{StatusCode: http.StatusBadRequest})
	}
}

func (r *RequestMock) GetUser() *models.User {
	return r.User
}

func (r *RequestMock) GetBodyMarshalled(ifc interface{}) {
	if err := json.Unmarshal(r.GetBody(), ifc); err != nil {
		panic(server.Error{StatusCode: http.StatusBadRequest})
	}
}

func (r *RequestMock) GetPathParameters() map[string]string {
	return r.PathParameters
}

func (r *RequestMock) GetParam(param string) string {
	return r.PathParameters[param]
}

func (r *RequestMock) GetParamUInt(param string) uint {
	if value, err := strconv.Atoi(r.PathParameters[param]); err == nil {
		return uint(value)
	} else {
		panic(err)
	}
}

func CreateFakeUser() *models.User {
	return &models.User{
		ID:        10,
		Name:      "name",
		Password:  []byte("hashedPassword"),
		Email:     "email@test.test",
		Active:    true,
		CreatedAt: null.Time{},
		UpdatedAt: null.Time{},
		DeletedAt: null.Time{},
	}
}
