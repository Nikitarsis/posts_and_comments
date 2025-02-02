package translationdao

type PostDao struct {
	PostId  uint64  `json:"post_id"`
	UserId  uint64  `json:"user_id"`
	Message *string `json:"message,omitempty"`
}
