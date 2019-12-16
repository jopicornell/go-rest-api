package requests

import "github.com/jopicornell/go-rest-api/db/entities/palmaactiva/image_gallery/model"

type EditUser struct {
	Username string `json:"username" validate:"required,alphanum,min=4"`
	FullName string `json:"fullname" validate:"ascii,min=5,required"`
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
