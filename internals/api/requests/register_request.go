package requests

import (
	"crypto/rand"
	"github.com/jopicornell/go-rest-api/internals/models"
	"golang.org/x/crypto/argon2"
	"gopkg.in/guregu/null.v3"
	"time"
)

type RegisterRequest struct {
	Name     string `json:"name" valid:"matches(^[a-zA-Z\\s]+$)"`
	Email    string `json:"email" valid:"email,required"`
	Password string `json:"password" valid:"required,gte=8,lte=72"`
}

func (r *RegisterRequest) TransformToUser() (user *models.User) {

	user = &models.User{
		ID:        0,
		Name:      r.Name,
		Password:  nil,
		Email:     r.Email,
		Active:    true,
		CreatedAt: null.TimeFrom(time.Now()),
		UpdatedAt: null.TimeFrom(time.Now()),
	}
	user.Password = generatePassword(r.Password)
	return
}

func generatePassword(password string) []byte {

	salt := generateRandomBytes(16)
	return argon2.IDKey([]byte(password), salt, 3, 3*1024, 2, 64)
}

func generateRandomBytes(n uint32) []byte {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return b
}
