package pictures

import (
	"fmt"
	errors2 "github.com/jopicornell/go-rest-api/api/pictures/errors"
	"github.com/jopicornell/go-rest-api/api/pictures/requests"
	"github.com/jopicornell/go-rest-api/api/pictures/services"
	"github.com/jopicornell/go-rest-api/api/users/errors"
	"github.com/jopicornell/go-rest-api/api/users/middlewares"
	"github.com/jopicornell/go-rest-api/api/users/models"
	"github.com/jopicornell/go-rest-api/pkg/server"
	"github.com/jopicornell/go-rest-api/pkg/utilities"
	"net/http"
)

type PictureHandler struct {
	server         server.Server
	pictureService services.PicturesService
}

func (a *PictureHandler) Initialize(server server.Server) {
	a.server = server
	a.pictureService = services.NewPictureService(server)
}

func (a *PictureHandler) ConfigureRoutes(router server.Router) {
	picturesGroup := router.AddGroup("/pictures")
	picturesGroup.AddRoute("", a.GetPictures).Methods(http.MethodGet)
	picturesGroup.AddRoute("/{id:[0-9]+}", a.GetOnePicture).Methods(http.MethodGet)
	picturesGroup.AddRoute("/{id:[0-9]+}/comments", a.GetPictureComments).Methods(http.MethodGet)
	picturesGroup.AddRoute("/{id:[0-9]+}/likes", a.GetPictureLikes).Methods(http.MethodGet)

	picturesRestrictedUserAdmin := picturesGroup.AddGroup("").
		Use(&middlewares.UserMiddleware{Roles: []string{models.USER_ROLE, models.ADMIN_ROLE}})
	picturesRestrictedUserAdmin.
		AddRoute("/{id:[0-9]+}", a.DeletePicture).
		Methods(http.MethodDelete)
	picturesRestrictedUserAdmin.
		AddRoute("/{id:[0-9]+}/comments/{comment_id:[0-9]+}", a.DeletePictureComment).
		Methods(http.MethodDelete)
	picturesRestrictedUserAdmin.
		AddRoute("/{id:[0-9]+}", a.UpdatePicture).
		Methods(http.MethodPut)

	picturesRestrictedUser := picturesGroup.AddGroup("").
		Use(&middlewares.UserMiddleware{Roles: []string{models.USER_ROLE}})
	picturesRestrictedUser.AddRoute("", a.CreatePicture).Methods(http.MethodPost)
	picturesRestrictedUser.AddRoute("/{id:[0-9]+}/comments", a.CreatePictureComment).Methods(http.MethodPost)
	picturesRestrictedUser.AddRoute("/{id:[0-9]+}/likes/{user_id:[0-9]+}", a.CreatePictureLike).Methods(http.MethodPost)

	picturesRestrictedUser.
		AddRoute("/{id:[0-9]+}/likes/{user_id:[0-9]+}", a.DeletePictureLike).
		Methods(http.MethodDelete)

}

func (a *PictureHandler) GetPictures(response server.Response, request server.Context) {
	if pictures, err := a.pictureService.GetPictures(); err == nil {
		response.RespondJSON(http.StatusOK, pictures)
	} else {
		response.Error(&server.Error{StatusCode: http.StatusInternalServerError, Error: err})
	}
}

func (a *PictureHandler) GetOnePicture(response server.Response, request server.Context) {
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

func (a *PictureHandler) UpdatePicture(response server.Response, request server.Context) {
	id := request.GetParamUInt("id")
	var updatePicture *requests.CreateOrUpdatePicture
	request.GetBodyMarshalled(&updatePicture)
	user := request.GetUser().(*models.UserWithRoles)
	if picture, err := a.pictureService.GetPicture(id); err != nil {
		switch err {
		case errors2.PictureNotFound:
			response.Respond(http.StatusNotFound)
		default:
			response.Error(&server.Error{StatusCode: http.StatusInternalServerError, Error: err})
		}
	} else if !a.pictureService.CheckUserAccess(user, picture) {
		response.Respond(http.StatusForbidden)
	} else if err := utilities.ValidateStruct(updatePicture); err != nil {
		response.RespondValidationErrors(http.StatusBadRequest, err)
	} else if picture, err := a.pictureService.UpdatePicture(uint(id), updatePicture.TransformToPicture()); err == nil {
		if picture == nil {
			response.Respond(http.StatusNotFound)
			return
		}
		response.RespondJSON(http.StatusOK, picture)
	} else {
		response.Error(&server.Error{StatusCode: http.StatusInternalServerError})
	}

}

func (a *PictureHandler) CreatePicture(response server.Response, request server.Context) {
	var createPictureRequest *requests.CreateOrUpdatePicture
	user := request.GetUser().(*models.UserWithRoles)
	request.GetBodyMarshalled(&createPictureRequest)
	if err := utilities.ValidateStruct(createPictureRequest); err != nil {
		response.RespondValidationErrors(http.StatusBadRequest, err)
	} else if picture, err := a.pictureService.CreatePicture(createPictureRequest.TransformToPicture(), user); err != nil {
		response.Error(&server.Error{StatusCode: http.StatusInternalServerError, Error: err})
	} else {
		response.SetHeader("Location", fmt.Sprintf("/pictures/%d", picture.PictureID))
		response.RespondJSON(http.StatusCreated, picture)
	}

}

func (a *PictureHandler) DeletePicture(response server.Response, request server.Context) {
	id := request.GetParamUInt("id")
	user := request.GetUser().(*models.UserWithRoles)
	if picture, err := a.pictureService.GetPicture(id); err != nil {
		switch err {
		case errors2.PictureNotFound:
			response.Respond(http.StatusNotFound)
		default:
			response.Error(&server.Error{StatusCode: http.StatusInternalServerError, Error: err})
		}
	} else if !a.pictureService.CheckUserAccess(user, picture) {
		response.Respond(http.StatusForbidden)
	} else if err := a.pictureService.DeletePicture(id, user); err != nil {
		switch err {
		case errors.ForbiddenAction:
			response.Respond(http.StatusForbidden)
		case errors2.PictureNotFound:
			response.Respond(http.StatusNotFound)
		default:
			response.Error(&server.Error{StatusCode: http.StatusInternalServerError, Error: err})
		}
	} else {
		response.Respond(http.StatusNoContent)
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
	id := request.GetParamUInt("id")
	user := request.GetUser().(*models.UserWithRoles)
	if picture, err := a.pictureService.GetPicture(id); err != nil {
		switch err {
		case errors2.PictureNotFound:
			response.Respond(http.StatusNotFound)
		default:
			response.Error(&server.Error{StatusCode: http.StatusInternalServerError, Error: err})
		}
	} else if !a.pictureService.CheckUserAccess(user, picture) {
		response.Respond(http.StatusForbidden)
	} else if err := a.pictureService.CreatePictureLike(int32(id), user.UserID); err != nil {
		switch err {
		case errors2.UserAlreadyLikedPicture:
			response.Respond(http.StatusConflict)
		default:
			response.Error(&server.Error{StatusCode: http.StatusInternalServerError, Error: err})
		}

	} else {
		response.SetHeader("Location", fmt.Sprintf("/pictures/%d/likes/%d", id, user.UserID))
		response.Respond(http.StatusCreated)
	}
}

func (a *PictureHandler) DeletePictureLike(response server.Response, request server.Context) {
	id := request.GetParamUInt("id")
	userId := request.GetParamInt("user_id")
	user := request.GetUser().(*models.UserWithRoles)
	if picture, err := a.pictureService.GetPicture(id); err != nil {
		switch err {
		case errors2.PictureNotFound:
			response.Respond(http.StatusNotFound)
		default:
			response.Error(&server.Error{StatusCode: http.StatusInternalServerError, Error: err})
		}
	} else if !a.pictureService.CheckUserAccess(user, picture) {
		response.Respond(http.StatusForbidden)
	} else if err := a.pictureService.DeletePictureLike(int32(id), int32(userId)); err != nil {
		switch err {
		case errors2.PictureNotFound:
			response.Respond(http.StatusNotFound)
		default:
			response.Error(&server.Error{StatusCode: http.StatusInternalServerError, Error: err})
		}
	} else {
		response.Respond(http.StatusNoContent)
	}
}

func (a *PictureHandler) DeletePictureComment(response server.Response, request server.Context) {
	pictureId := request.GetParamInt("id")
	commentId := request.GetParamInt("comment_id")
	user := request.GetUser().(*models.UserWithRoles)
	if picture, err := a.pictureService.GetPicture(uint(pictureId)); err != nil {
		switch err {
		case errors2.PictureNotFound:
			response.Respond(http.StatusNotFound)
		default:
			response.Error(&server.Error{StatusCode: http.StatusInternalServerError, Error: err})
		}
	} else if !a.pictureService.CheckUserAccess(user, picture) {
		response.Respond(http.StatusForbidden)
	} else if err := a.pictureService.DeletePictureComment(int32(pictureId), int32(commentId)); err == nil {
		response.Respond(http.StatusNoContent)
	} else {
		switch err {
		case errors2.PictureNotFound:
			response.Respond(http.StatusNotFound)
		default:
			response.Error(&server.Error{StatusCode: http.StatusInternalServerError, Error: err})
		}
	}
}
