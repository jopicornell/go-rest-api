package responses

type Comment struct {
	Comment   string `json:"comment"`
	PictureID string `json:"picture_id"`
	UserID    string `json:"user_id"`
}
