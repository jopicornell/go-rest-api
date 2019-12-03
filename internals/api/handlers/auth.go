package handlers

import (
	"github.com/asaskevich/govalidator"
	"github.com/jopicornell/go-rest-api/internals/api/requests"
	"github.com/jopicornell/go-rest-api/internals/api/services"
	"github.com/jopicornell/go-rest-api/internals/errors"
	"github.com/jopicornell/go-rest-api/pkg/server"
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
			response.Respond(http.StatusUnauthorized)
		default:
			response.Error(&server.Error{StatusCode: http.StatusInternalServerError, Error: err})
		}
	}
}

func (a *AuthHandler) ConfigureRoutes(router server.Router) {
	router.AddRoute("/login", a.Login).Methods("POST")
}
