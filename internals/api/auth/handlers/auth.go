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
	server.Handler
	authService services.AuthService
}

func New(s server.Server) *AuthHandler {
	return &AuthHandler{
		authService: services.New(s.GetRelationalDatabase(), s),
	}
}

func (a *AuthHandler) Login(context server.Context) {
	var loginRequest requests.LoginRequest
	context.GetBodyMarshalled(&loginRequest)
	valid, err := govalidator.ValidateStruct(loginRequest)
	if err != nil || !valid {
		context.Respond(http.StatusBadRequest)
		return
	}
	if token, err := a.authService.Login(loginRequest.Email, loginRequest.Password); err == nil {
		context.RespondJSON(http.StatusOK, token)
	} else {
		switch err {
		case errors.AuthUserNotMatched:
			context.Respond(http.StatusForbidden)
		default:
			context.Respond(http.StatusInternalServerError)
		}
	}
}
