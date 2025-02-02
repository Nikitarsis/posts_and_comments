package translationdao

type CommentDao struct {
	UserId      uint64   `json:"user_id"`
	PostId      uint64   `json:"post_id"`
	ParentId    uint64   `json:"parent_id"`
	ChildrenIds []uint64 `json:"children_id"`
	Message     *string  `json:"message,omitempty"`
}
