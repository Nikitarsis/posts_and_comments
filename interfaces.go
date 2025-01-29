package comments_and_posts

type usrId uint64
type msgId uint64
type content string

type IHaveMessageId interface {
	GetMessageId() msgId
}

type IHaveContent interface {
	GetContent() content
	SetContent(data content)
}

type IHaveComments interface {
	GetCommentsId() ([]msgId, error)
	AddCommentId(childId ...msgId) error
	RemoveCommentId(id ...msgId) error
}

type IHaveParent interface {
	GetParentId() (msgId, error)
	SetParentId(id msgId) error
}

type IHaveChildren interface {
	GetChildrenId() (msgId, error)
	SetChildrenId(id msgId) error
}

type IMessage interface {
	IHaveMessageId
	IHaveContent
}

type IPost interface {
	IHaveMessageId
	IHaveComments
}

type IComment interface {
	IHaveParent
	IHaveChildren
}
