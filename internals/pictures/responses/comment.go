package responses

type Comment struct {
	CommentID int32  `json:"comment_id"`
	Comment   string `json:"comment"`
	PictureID int32  `json:"picture_id"`
	UserID    int32  `json:"user_id"`
}
