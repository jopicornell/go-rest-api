package requests

import (
	"github.com/jopicornell/go-rest-api/db/entities/palmaactiva/image_gallery/model"
	password "github.com/jopicornell/go-rest-api/pkg/password"
	"github.com/sirupsen/logrus"
)

type RegisterRequest struct {
	FullName string `json:"fullname" valid:"matches(^[a-zA-Z\\s]+$)"`
	Username string `json:"username" valid:"alphanum,required"`
	Password string `json:"password" valid:"required,gte=6,lte=72"`
}

func (r *RegisterRequest) TransformToUser() (user *model.User) {

	user = &model.User{
		FullName: r.FullName,
		Password: generatePassword(r.Password),
		Username: r.Username,
	}
	return
}

func generatePassword(passwd string) string {
	params := &password.ArgonPasswordParams{
		Memory:      64 * 1024,
		Iterations:  3,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	}
	if hashedPassword, err := password.ArgonHashFromPassword(passwd, params); err != nil {
		logrus.Panic("Error hashing password: %w", err)
	} else {
		return hashedPassword
	}
	return ""
}
