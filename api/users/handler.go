package users

import (
	"github.com/jopicornell/go-rest-api/api/users/errors"
	"github.com/jopicornell/go-rest-api/api/users/middlewares"
	"github.com/jopicornell/go-rest-api/api/users/models"
	"github.com/jopicornell/go-rest-api/api/users/requests"
	"github.com/jopicornell/go-rest-api/api/users/services"
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
	a.authService = services.NewAuthService(server)
}

func (a *UserHandler) ConfigureRoutes(router server.Router) {
	userGroup := router.AddGroup("/users")
	userGroup.AddRoute("/login", a.Login).Methods(http.MethodPost)
	userGroup.AddRoute("", a.Register).Methods(http.MethodPost)
	restrictedUserAdmin := userGroup.AddGroup("").
		Use(&middlewares.UserMiddleware{Roles: []string{models.USER_ROLE, models.ADMIN_ROLE}})
	restrictedUserAdmin.AddRoute("/{id:[0-9]+}", a.GetUser).Methods(http.MethodGet)
	restrictedUserAdmin.AddRoute("/{id:[0-9]+}", a.UpdateUser).Methods(http.MethodPut)
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

func (a *UserHandler) GetUser(response server.Response, context server.Context) {
	id := context.GetParamUInt("id")
	if picture, err := a.authService.GetUser(id); err == nil {
		if picture == nil {
			response.Respond(http.StatusNotFound)
			return
		}
		response.RespondJSON(http.StatusOK, picture)
	} else {
		response.Error(&server.Error{StatusCode: http.StatusInternalServerError, Error: err})
	}
}

func (a *UserHandler) UpdateUser(response server.Response, context server.Context) {
	id := context.GetParamUInt("id")
	var editUserRequest requests.EditUser
	context.GetBodyMarshalled(&editUserRequest)
	if picture, err := a.authService.UpdateUser(id, editUserRequest.TransformToUser()); err == nil {
		if picture == nil {
			response.Respond(http.StatusNotFound)
			return
		}
		response.RespondJSON(http.StatusOK, picture)
	} else {
		response.Error(&server.Error{StatusCode: http.StatusInternalServerError, Error: err})
	}
}
