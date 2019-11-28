package requests

type LoginRequest struct {
	Email    string `json:"email" valid:"email,required"`
	Password string `json:"password" valid:"required"`
}
