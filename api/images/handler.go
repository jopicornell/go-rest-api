package images

import (
	"fmt"
	"github.com/jopicornell/go-rest-api/api/images/services"
	"github.com/jopicornell/go-rest-api/api/users/middlewares"
	"github.com/jopicornell/go-rest-api/api/users/models"
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

func (ih *ImageHandler) SaveImage(res server.Response, context server.Context) {
	imageType := context.GetRequest().URL.Query().Get("type")
	user := context.GetUser().(*models.UserWithRoles)
	image := ih.imageService.SavePicture(user, context.GetBody(), imageType)
	res.SetHeader("Location", fmt.Sprintf("/images/%d", image.ImageID))
	res.RespondJSON(http.StatusCreated, image)
}
