package responses

type User struct {
	UserID   int32  `json:"user_id"`
	Username string `json:"username"`
	FullName string `json:"full_name"`
}
