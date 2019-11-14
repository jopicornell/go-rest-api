package servertesting

import (
	"encoding/json"
	"github.com/jopicornell/go-rest-api/internals/models"
	"gopkg.in/guregu/null.v3"
	"io/ioutil"
	"net/http"
	"strconv"
)

func NewContext(r *http.Request, w http.ResponseWriter, pathParameters map[string]string) *ContextMock {
	return &ContextMock{
		Request:        r,
		ResponseWriter: w,
		PathParameters: pathParameters,
	}
}

type ContextMock struct {
	Request        *http.Request
	PathParameters map[string]string
	ResponseWriter http.ResponseWriter
	ThrowError     error
	User           *models.User
}

func (r *ContextMock) GetRequest() *http.Request {
	return r.Request
}

func (r *ContextMock) GetWriter() http.ResponseWriter {
	return r.ResponseWriter
}

func (r *ContextMock) GetBody() []byte {
	if body, err := ioutil.ReadAll(r.Request.Body); err == nil {
		return body
	} else {
		r.Respond(http.StatusBadRequest)
		return nil
	}
}

func (r *ContextMock) GetUser() *models.User {
	return r.User
}

func (r *ContextMock) GetBodyMarshalled(ifc interface{}) {
	if err := json.Unmarshal(r.GetBody(), ifc); err != nil {
		r.Respond(http.StatusBadRequest)
	}
}

func (r *ContextMock) GetPathParameters() map[string]string {
	return r.PathParameters
}

func (r *ContextMock) GetParam(param string) string {
	return r.PathParameters[param]
}

func (r *ContextMock) GetParamUInt(param string) uint {
	if value, err := strconv.Atoi(r.PathParameters[param]); err == nil {
		return uint(value)
	} else {
		panic(err)
	}
}

func (r *ContextMock) RespondJSON(statusCode int, ifc interface{}) {
	r.ResponseWriter.Header().Add("Content-Type", "application/json")
	r.ResponseWriter.WriteHeader(statusCode)
	if err := json.NewEncoder(r.ResponseWriter).Encode(ifc); err == nil {
	} else {
		panic(err)
	}
}

func (r *ContextMock) Respond(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
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
