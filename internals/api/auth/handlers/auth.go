package handlers

import (
	"github.com/asaskevich/govalidator"
	"github.com/jopicornell/go-rest-api/internals/api/auth/requests"
	"github.com/jopicornell/go-rest-api/internals/api/auth/services"
	"github.com/jopicornell/go-rest-api/internals/errors"
	"github.com/jopicornell/go-rest-api/pkg/server"
	"net/http"
)

type AuthHandler struct {
	authService services.AuthService
}

func New(s server.Server) *AuthHandler {
	return &AuthHandler{
		authService: services.New(s.GetRelationalDatabase(), s),
	}
}

func (a *AuthHandler) Login(response server.Response, request server.Request) {
	var loginRequest requests.LoginRequest
	request.GetBodyMarshalled(&loginRequest)
	valid, err := govalidator.ValidateStruct(loginRequest)
	if err != nil || !valid {
		response.Respond(http.StatusBadRequest)
		return
	}
	if token, err := a.authService.Login(loginRequest.Email, loginRequest.Password); err == nil {
		response.RespondJSON(http.StatusOK, token)
	} else {
		switch err {
		case errors.AuthUserNotMatched:
			response.Respond(http.StatusForbidden)
		default:
			response.Error(&server.Error{StatusCode: http.StatusInternalServerError, Error: err})
		}
	}
}
