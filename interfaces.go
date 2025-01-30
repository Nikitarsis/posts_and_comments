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
	GetCommentsId() []msgId
	AddCommentsId(ids ...msgId)
	RemoveCommentsId(ids ...msgId)
	GetCommentsNum() uint
}

type IHaveParent interface {
	GetParentId() msgId
	SetParentId(id msgId)
}

type IHaveChildren interface {
	GetChildrenIds() []msgId
	AddChildrenIds(ids ...msgId)
	RemoveChildrenIds(ids ...msgId)
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
	IHaveMessageId
	IHaveParent
	IHaveChildren
}
