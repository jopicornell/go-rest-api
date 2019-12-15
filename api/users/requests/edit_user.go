package requests

import "github.com/jopicornell/go-rest-api/db/entities/palmaactiva/image_gallery/model"

type EditUser struct {
	Username string `json:"username" validate:"required"`
	FullName string `json:"fullname" validate:"required"`
	AvatarId *int32 `json:"image_id" validate:"required,numeric"`
}

func (r *EditUser) TransformToUser() (user *model.User) {

	user = &model.User{
		FullName: r.FullName,
		AvatarID: r.AvatarId,
		Username: r.Username,
	}
	return
}
