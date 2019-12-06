package handlers

import (
	"github.com/jopicornell/go-rest-api/internals/api/middlewares"
	"github.com/jopicornell/go-rest-api/internals/api/services"
	"github.com/jopicornell/go-rest-api/internals/models"
	"github.com/jopicornell/go-rest-api/pkg/server"
	"net/http"
)

type PictureHandler struct {
	server         server.Server
	pictureService services.PicturesService
}

func (a *PictureHandler) Initialize(server server.Server) {
	a.server = server
	a.pictureService = services.NewPictureService(server.GetRelationalDatabase())
}

func (a *PictureHandler) ConfigureRoutes(router server.Router) {
	pictures := router.AddGroup("/pictures")
	pictures.Use(&middlewares.UserMiddleware{}, &middlewares.UserMiddleware{})
	pictures.AddRoute("", a.GetPicturesHandler).Methods("GET")
	pictures.AddRoute("", a.CreatePictureHandler).Methods("POST")
	pictures.AddRoute("/{id:[0-9]+}", a.DeletePictureHandler).Methods("DELETE")
	pictures.AddRoute("/{id:[0-9]+}", a.GetOnePictureHandler).Methods("GET")
	pictures.AddRoute("/{id:[0-9]+}", a.UpdatePictureHandler).Methods("PUT")
}

func (a *PictureHandler) GetPicturesHandler(response server.Response, request server.Context) {
	if pictures, err := a.pictureService.GetPictures(); err == nil {
		response.RespondJSON(http.StatusOK, pictures)
	} else {
		response.Error(&server.Error{StatusCode: http.StatusInternalServerError, Error: err})
	}
}

func (a *PictureHandler) GetOnePictureHandler(response server.Response, request server.Context) {
	id := request.GetParamUInt("id")
	if picture, err := a.pictureService.GetPicture(uint(id)); err == nil {
		if picture == nil {
			response.Respond(http.StatusNotFound)
			return
		}
		response.RespondJSON(http.StatusOK, picture)
	} else {
		response.Error(&server.Error{StatusCode: http.StatusInternalServerError, Error: err})
	}
}

func (a *PictureHandler) UpdatePictureHandler(response server.Response, request server.Context) {
	id := request.GetParamUInt("id")
	var picture *models.Picture
	request.GetBodyMarshalled(&picture)
	if picture, err := a.pictureService.UpdatePicture(uint(id), picture); err == nil {
		if picture == nil {
			response.Respond(http.StatusNotFound)
			return
		}
		response.RespondJSON(http.StatusOK, picture)
	} else {
		response.Error(&server.Error{StatusCode: http.StatusInternalServerError})
	}

}

func (a *PictureHandler) CreatePictureHandler(response server.Response, request server.Context) {
	var picture *models.Picture
	request.GetBodyMarshalled(&picture)
	if picture, err := a.pictureService.CreatePicture(picture, request.GetUser()); err == nil {
		response.RespondJSON(http.StatusCreated, picture)
	} else {
		response.Error(&server.Error{StatusCode: http.StatusInternalServerError, Error: err})
	}

}

func (a *PictureHandler) DeletePictureHandler(response server.Response, request server.Context) {
	id := request.GetParamUInt("id")
	if err := a.pictureService.DeletePicture(uint(id)); err == nil {
		response.Respond(http.StatusNoContent)
	} else {
		response.Error(&server.Error{StatusCode: http.StatusInternalServerError, Error: err})
	}
}
