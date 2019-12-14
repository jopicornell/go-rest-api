package users

import (
	"github.com/jopicornell/go-rest-api/internals/errors"
	"github.com/jopicornell/go-rest-api/internals/users/requests"
	"github.com/jopicornell/go-rest-api/internals/users/services"
	"github.com/jopicornell/go-rest-api/pkg/server"
	"github.com/jopicornell/go-rest-api/pkg/utilities"
	"github.com/sirupsen/logrus"
	"net/http"
)

type UserHandler struct {
	server      server.Server
	authService services.UsersService
}

func (a *UserHandler) Initialize(server server.Server) {
	a.server = server
	a.authService = services.NewAuthService(server.GetRelationalDatabase(), server)
}

func (a *UserHandler) ConfigureRoutes(router server.Router) {
	userGroup := router.AddGroup("/users")
	userGroup.AddRoute("/login", a.Login).Methods("POST")
	userGroup.AddRoute("", a.Register).Methods("POST")
}

func (a *UserHandler) Login(response server.Response, request server.Context) {
	var loginRequest requests.LoginRequest
	request.GetBodyMarshalled(&loginRequest)
	err := utilities.ValidateStruct(loginRequest)
	if err != nil {
		response.RespondValidationErrors(http.StatusBadRequest, err)
		return
	}
	if token, err := a.authService.Login(loginRequest.Username, loginRequest.Password); err == nil {
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

func (a *UserHandler) Register(response server.Response, context server.Context) {
	var registerRequest requests.RegisterRequest
	context.GetBodyMarshalled(&registerRequest)
	err := utilities.ValidateStruct(registerRequest)
	if err != nil {
		response.RespondJSON(http.StatusBadRequest, err)
	}
	if user, err := a.authService.Register(registerRequest.TransformToUser()); err == nil {
		response.RespondJSON(http.StatusCreated, user)
	} else {
		if err == errors.UsernameExists {
			response.Respond(409)
		} else {
			logrus.Error(err)
			response.Respond(http.StatusInternalServerError)
		}
	}
}
