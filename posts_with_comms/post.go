package comments_and_posts

import (
	"reflect"

	"../messages"
)

/*
Простая реализация IPost, указывающая на зависимости между сообщениями
Вместо ссылок на другие объекты Post используется не нарушающая
принцип единой ответственности схема с ID
*/
type Post struct {
	id       messages.MsgId
	parent   messages.MsgId //ID родителя, если указан id поста, родителя нет
	children map[messages.MsgId]struct{}
}

/*
Возвращает ID сообщения поста
*/
func (c Post) GetMessageId() messages.MsgId {
	return c.id
}

/*
Возвращает ID сообщения родителя и false, либо ID сообщения поста и true, если родителя нет
*/
func (c Post) GetParentId() (messages.MsgId, bool) {
	return c.parent, reflect.DeepEqual(c.id, c.parent)
}

/*
Получает дочерние ID
*/
func (c Post) GetChildrenIds() []messages.MsgId {
	ret := make([]messages.MsgId, len(c.children))
	i := 0
	for k := range c.children {
		ret[i] = k
		i++
	}
	return ret
}

/*
Добавляет дочерние ID
*/
func (c *Post) AddChildrenIds(ids ...messages.MsgId) {
	for _, id := range ids {
		c.children[id] = struct{}{}
	}
}

/*
Простейший конструктор начального поста без родителя
*/
func NewInitPost(id messages.MsgId) *Post {
	return &Post{id, id, make(map[messages.MsgId]struct{})}
}

/*
Конструктор поста с родителем
*/
func NewPost(id messages.MsgId, parent messages.MsgId) *Post {
	return &Post{id, parent, make(map[messages.MsgId]struct{})}
}

/*
Конструктор поста с родителем и дочерними элементами
*/
func NewPostWithChildren(id messages.MsgId, parent messages.MsgId, children ...messages.MsgId) *Post {
	childrenMap := make(map[messages.MsgId]struct{}, len(children))
	for _, child := range children {
		childrenMap[child] = struct{}{}
	}
	return &Post{id, parent, childrenMap}
}
