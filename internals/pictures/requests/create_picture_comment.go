package requests

import "github.com/jopicornell/go-rest-api/db/entities/palmaactiva/image_gallery/model"

type CreatePictureComment struct {
	Comment   string `json:"comment"`
	PictureID int32  `json:"picture_id"`
	UserID    int32  `json:"user_id"`
}

func (cpc *CreatePictureComment) TransformToComment() *model.Comment {
	return &model.Comment{
		UserID:    cpc.UserID,
		PictureID: cpc.PictureID,
		Comment:   cpc.Comment,
	}
}
