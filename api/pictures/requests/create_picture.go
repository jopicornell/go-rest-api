package requests

import "github.com/jopicornell/go-rest-api/db/entities/palmaactiva/image_gallery/model"

type CreateOrUpdatePicture struct {
	Description string `json:"description"`
	ImageID     int    `json:"image_id" validate:"numeric"`
	Title       string `json:"title" validate:"min=3"`
	UserID      int    `json:"user_id" validate:"numeric"`
}

func (cpr *CreateOrUpdatePicture) TransformToPicture() *model.Picture {
	return &model.Picture{
		UserID:      int32(cpr.UserID),
		ImageID:     int32(cpr.ImageID),
		Title:       cpr.Title,
		Description: &cpr.Description,
	}
}
