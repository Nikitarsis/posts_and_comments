package comments_and_posts

/*
Основные интерфейсы пакета
*/

// Обложка для ID
type msgId uint64

// Обложка для сообщения
type content string

// Структура имеет ID сообщения
type IHaveMessageId interface {
	GetMessageId() msgId
}

// Структура содержит содержание(текст, ссылки на прикреплёные файлы etc.)
type IHaveContent interface {
	GetContent() content
	SetContent(data content)
}

// Структура поддерживает комментарии
type IHaveComments interface {
	GetCommentsId() []msgId
	AddCommentsId(ids ...msgId)
	RemoveCommentsId(ids ...msgId)
	GetCommentsNum() uint
}

// Структура может иметь родительский элемент(для комментариев-ответов)
type IHaveParent interface {
	GetParentId() (msgId, bool)
	SetParentId(id msgId)
}

// Структура может иметь дочерние элементы(для комментариев с ответами)
type IHaveChildren interface {
	GetChildrenIds() []msgId
	AddChildrenIds(ids ...msgId)
	RemoveChildrenIds(ids ...msgId)
}

// Сообщение с Id и содержанием
type IMessage interface {
	IHaveMessageId
	IHaveContent
}

// Пост, который может иметь комментарии
type IPost interface {
	IHaveMessageId
	IHaveComments
}

// Комментарий, который может ссылаться
type IComment interface {
	IHaveMessageId
	IHaveParent
	IHaveChildren
}
