package comments_and_posts

import "github.com/Nikitarsis/posts_and_comments/messages"

/*
Основные интерфейсы пакета
*/

// Обложка для сообщения
type content string

// Структура имеет ID сообщения
type IHaveMessageId interface {
	//Возвращает свой ID
	GetMessageId() messages.MsgId
}

// Структура содержит содержание(текст, ссылки на прикреплёные файлы etc.)
type IHaveContent interface {
	//Возвращает содержание
	GetContent() content
	//Задаёт содержание
	SetContent(data content)
}

// Структура может иметь родительский элемент(для комментариев-ответов)
type IHaveParent interface {
	//Возвращает ID родителя и false, либо другой ID(напр. собственный) и true, если родителя нет
	GetParentId() (messages.MsgId, bool)
}

// Структура может иметь дочерние элементы(для комментариев с ответами)
type IHaveChildren interface {
	//Возвращает ID дочерних собщений
	GetChildrenIds() []messages.MsgId
	//Добавляет дочерние сообщения
	AddChildrenIds(ids ...messages.MsgId)
}

// Сообщение с Id и содержанием
type IMessage interface {
	IHaveMessageId
	IHaveContent
}

// Комментарий, который может ссылаться
type IPost interface {
	IHaveMessageId
	IHaveParent
	IHaveChildren
}
