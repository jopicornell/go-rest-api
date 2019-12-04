package handlers

import (
	"github.com/jopicornell/go-rest-api/internals/api/requests"
	"github.com/jopicornell/go-rest-api/internals/api/services"
	"github.com/jopicornell/go-rest-api/internals/errors"
	"github.com/jopicornell/go-rest-api/pkg/server"
	"github.com/jopicornell/go-rest-api/pkg/utilities"
	"github.com/sirupsen/logrus"
	"net/http"
)

type AuthHandler struct {
	server      server.Server
	authService services.AuthService
}

func (a *AuthHandler) Initialize(server server.Server) {
	a.server = server
	a.authService = services.NewAuthService(server.GetRelationalDatabase(), server)
}

func (a *AuthHandler) ConfigureRoutes(router server.Router) {
	router.AddRoute("/login", a.Login).Methods("POST")
	router.AddRoute("/register", a.Register).Methods("POST")
}

func (a *AuthHandler) Login(response server.Response, request server.Context) {
	var loginRequest requests.LoginRequest
	request.GetBodyMarshalled(&loginRequest)
	err := utilities.ValidateStruct(loginRequest)
	if err != nil {
		response.RespondValidationErrors(http.StatusBadRequest, err)
		return
	}
	if token, err := a.authService.Login(loginRequest.Email, loginRequest.Password); err == nil {
		response.RespondJSON(http.StatusOK, token)
	} else {
		switch err {
		case errors.AuthUserNotMatched:
			response.Respond(http.StatusUnauthorized)
		default:
			response.Error(&server.Error{StatusCode: http.StatusInternalServerError, Error: err})
		}
	}
}

func (a *AuthHandler) Register(response server.Response, context server.Context) {
	var registerRequest requests.RegisterRequest
	context.GetBodyMarshalled(&registerRequest)
	err := utilities.ValidateStruct(registerRequest)
	if err != nil {
		response.RespondJSON(http.StatusBadRequest, err)
	}
	if user, err := a.authService.Register(registerRequest.TransformToUser()); err == nil {
		response.RespondJSON(http.StatusOK, user)
	} else {
		logrus.Error(err)
		response.Respond(http.StatusBadRequest)
	}
}
