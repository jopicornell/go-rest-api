package pictures

import (
	"fmt"
	"github.com/jopicornell/go-rest-api/db/entities/palmaactiva/image_gallery/model"
	"github.com/jopicornell/go-rest-api/internals/errors"
	"github.com/jopicornell/go-rest-api/internals/models"
	"github.com/jopicornell/go-rest-api/internals/pictures/requests"
	"github.com/jopicornell/go-rest-api/internals/pictures/services"
	"github.com/jopicornell/go-rest-api/internals/users/middlewares"
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
	picturesGroup := router.AddGroup("/pictures")
	picturesGroup.AddRoute("", a.GetPicturesHandler).Methods(http.MethodGet)
	picturesGroup.AddRoute("/{id:[0-9]+}", a.GetOnePictureHandler).Methods(http.MethodGet)
	picturesGroup.AddRoute("/{id:[0-9]+}/comments", a.GetPictureComments).Methods(http.MethodGet)
	picturesGroup.AddRoute("/{id:[0-9]+}/likes", a.GetPictureLikes).Methods(http.MethodGet)

	picturesRestrictedUserAdmin := picturesGroup.AddGroup("").
		Use(&middlewares.UserMiddleware{Roles: []string{models.USER_ROLE, models.ADMIN_ROLE}})
	picturesRestrictedUserAdmin.
		AddRoute("/{id:[0-9]+}", a.DeletePicture).
		Methods(http.MethodDelete)
	picturesRestrictedUserAdmin.
		AddRoute("/{id:[0-9]+}/likes/{user_id:[0-9]+}", a.DeletePictureLike).
		Methods(http.MethodDelete)

	picturesRestrictedUser := picturesGroup.AddGroup("").
		Use(&middlewares.UserMiddleware{Roles: []string{models.USER_ROLE}})
	picturesRestrictedUser.AddRoute("", a.CreatePictureHandler).Methods(http.MethodPost)
	picturesRestrictedUser.AddRoute("/{id:[0-9]+}/comments", a.CreatePictureComment).Methods(http.MethodPost)
	picturesRestrictedUser.AddRoute("/{id:[0-9]+}/likes/{user_id:[0-9]+}", a.CreatePictureLike).Methods(http.MethodPost)

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
	if picture, err := a.pictureService.GetPicture(id); err == nil {
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
	var picture *model.Picture
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
	var createPictureRequest *requests.CreatePicture
	user := request.GetUser().(*models.UserWithRoles)
	request.GetBodyMarshalled(&createPictureRequest)
	if picture, err := a.pictureService.CreatePicture(createPictureRequest.TransformToPicture(), user); err == nil {
		response.SetHeader("Location", fmt.Sprintf("/pictures/%d", picture.PictureID))
		response.RespondJSON(http.StatusCreated, picture)
	} else {
		response.Error(&server.Error{StatusCode: http.StatusInternalServerError, Error: err})
	}

}

func (a *PictureHandler) DeletePicture(response server.Response, request server.Context) {
	id := request.GetParamUInt("id")
	user := request.GetUser().(*models.UserWithRoles)
	if err := a.pictureService.DeletePicture(id, user); err == nil {
		response.Respond(http.StatusNoContent)
	} else {
		switch err {
		case errors.ForbiddenAction:
			response.Respond(http.StatusForbidden)
		case errors.PictureNotFound:
			response.Respond(http.StatusNotFound)
		default:
			response.Error(&server.Error{StatusCode: http.StatusInternalServerError, Error: err})
		}
	}
}

func (a *PictureHandler) GetPictureComments(response server.Response, request server.Context) {
	id := request.GetParamUInt("id")
	if picture, err := a.pictureService.GetPictureComments(id); err == nil {
		if picture == nil {
			response.Respond(http.StatusNotFound)
			return
		}
		response.RespondJSON(http.StatusOK, picture)
	} else {
		response.Error(&server.Error{StatusCode: http.StatusInternalServerError, Error: err})
	}
}

func (a *PictureHandler) CreatePictureComment(response server.Response, request server.Context) {
	id := request.GetParamInt("id")
	var createPictureComment *requests.CreatePictureComment
	request.GetBodyMarshalled(&createPictureComment)
	if comment, err := a.pictureService.CreatePictureComment(int32(id), createPictureComment.TransformToComment()); err == nil {
		response.SetHeader("Location", fmt.Sprintf("/pictures/%d/comments/%d", id, comment.CommentID))
		response.RespondJSON(http.StatusCreated, comment)
	} else {
		response.Error(&server.Error{StatusCode: http.StatusInternalServerError, Error: err})
	}
}

func (a *PictureHandler) GetPictureLikes(response server.Response, request server.Context) {
	id := request.GetParamUInt("id")
	if picture, err := a.pictureService.GetPictureLikes(id); err == nil {
		if picture == nil {
			response.Respond(http.StatusNotFound)
			return
		}
		response.RespondJSON(http.StatusOK, picture)
	} else {
		response.Error(&server.Error{StatusCode: http.StatusInternalServerError, Error: err})
	}
}

func (a *PictureHandler) CreatePictureLike(response server.Response, request server.Context) {
	id := request.GetParamInt("id")
	userId := request.GetParamInt("user_id")
	if comment, err := a.pictureService.CreatePictureLike(int32(id), int32(userId)); err == nil {
		response.SetHeader("Location", fmt.Sprintf("/pictures/%d/likes/%d", id, comment.UserID))
		response.RespondJSON(http.StatusCreated, comment)
	} else {
		switch err {
		case errors.UserAlreadyLikedPicture:
			response.Respond(http.StatusConflict)
		default:
			response.Error(&server.Error{StatusCode: http.StatusInternalServerError, Error: err})
		}
	}
}

func (a *PictureHandler) DeletePictureLike(response server.Response, request server.Context) {
	id := request.GetParamInt("id")
	userId := request.GetParamInt("user_id")
	if err := a.pictureService.DeletePictureLike(int32(id), int32(userId)); err == nil {
		response.Respond(http.StatusNoContent)
	} else {
		switch err {
		case errors.PictureNotFound:
			response.Respond(http.StatusNotFound)
		default:
			response.Error(&server.Error{StatusCode: http.StatusInternalServerError, Error: err})
		}
	}
}
