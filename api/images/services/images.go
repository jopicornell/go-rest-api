package services

import (
	"github.com/google/uuid"
	"github.com/graux/image-manager"
	"github.com/jopicornell/go-rest-api/api/pictures/responses"
	"github.com/jopicornell/go-rest-api/api/users/models"
	. "github.com/jopicornell/go-rest-api/db/entities/palmaactiva/image_gallery/table"
	"github.com/jopicornell/go-rest-api/pkg/server"
	"github.com/sirupsen/logrus"
	"path/filepath"
)

func NewImagesService(server server.Server) ImageService {
	return &imagesService{server: server}
}

type ImageService interface {
	SavePicture(user *models.UserWithRoles, imageBytes []byte, imageType string) *responses.Image
}

type imagesService struct {
	server server.Server
}

func (is *imagesService) SavePicture(user *models.UserWithRoles, imageBytes []byte, imageType string) *responses.Image {
	imagesPath, err := filepath.Abs("./images")
	if err != nil {
		logrus.Panicf("error getting the images folder: %s", err)
	}
	imageManager := image_manger.NewImageManager(imagesPath)
	var uuids []uuid.UUID
	if imageType == "avatar" {
		uuids, err = imageManager.ProcessImageAsSquare(imageBytes)
	} else {
		uuids, err = imageManager.ProcessImageAs16by9(imageBytes)
	}
	if err != nil {
		logrus.Panicf("error processing the images: %s", err)
	}
	thumb := "files/images/" + uuids[0].String() + ".jpg"
	lowRes := "files/images/" + uuids[1].String() + ".jpg"
	highRes := "files/images/" + uuids[2].String() + ".jpg"
	statement := Image.INSERT(Image.UserID, Image.LowResURL, Image.HighResURL, Image.ThumbURL).
		VALUES(user.UserID, lowRes, highRes, thumb).
		RETURNING(Image.AllColumns)
	image := new(responses.Image)
	if err := statement.Query(is.server.GetRelationalDatabase().DB, image); err != nil {
		logrus.Panicf("error when saving the images to db: %s", err)
	} else {
		return image
	}
	return nil
}
