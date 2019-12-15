package requests

import "github.com/jopicornell/go-rest-api/db/entities/palmaactiva/image_gallery/model"

type CreatePicture struct {
	Description string `json:"description"`
	Image       string `json:"image"`
	ImageID     int    `json:"image_id"`
	Title       string `json:"title"`
	UserID      int    `json:"user_id"`
}

func (cpr *CreatePicture) TransformToPicture() *model.Picture {
	return &model.Picture{
		UserID:      int32(cpr.UserID),
		ImageID:     int32(cpr.ImageID),
		Title:       cpr.Title,
		Description: &cpr.Description,
	}
}
