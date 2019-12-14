package images

import (
	"fmt"
	"github.com/jopicornell/go-rest-api/internals/images/services"
	"github.com/jopicornell/go-rest-api/internals/models"
	"github.com/jopicornell/go-rest-api/internals/users/middlewares"
	"github.com/jopicornell/go-rest-api/pkg/server"
	"net/http"
)

type ImageHandler struct {
	server.Handler
	server       server.Server
	imageService services.ImageService
}

func (ih *ImageHandler) Initialize(server server.Server) {
	ih.server = server
	ih.imageService = services.NewImagesService(server)
}

func (ih *ImageHandler) ConfigureRoutes(router server.Router) {
	imagesGroup := router.AddGroup("/images")
	imagesGroup.Use(&middlewares.UserMiddleware{
		Roles: []string{"user"},
	})
	imagesGroup.AddRoute("", ih.SaveImage).Methods(http.MethodPost)
}

func (ih *ImageHandler) SaveImage(res server.Response, ctx server.Context) {
	user := ctx.GetUser().(*models.UserWithRoles)
	image := ih.imageService.SavePicture(user, ctx.GetBody())
	res.SetHeader("Location", fmt.Sprintf("/images/%d", image.ImageID))
	res.RespondJSON(http.StatusCreated, image)
}
