package comments_and_posts

import "cmp"

/*
Простая реализация IPost, указывающая на зависимости между сообщениями
*/
type Post struct {
	id       msgId
	parent   msgId //ID родителя, если указан id поста, родителя нет
	children map[msgId]struct{}
}

/*
Возвращает ID сообщения поста
*/
func (c Post) GetMessageId() msgId {
	return c.id
}

/*
Возвращает ID сообщения родителя и false, либо ID сообщения поста и true, если родителя нет
*/
func (c Post) GetParentId() (msgId, bool) {
	return c.parent, cmp.Compare(c.id, c.parent) == 0
}

/*
Устанавливает ID сообщения родителя
*/
//TODO: Подумать над необходимостью
func (c *Post) SetParentId(id msgId) {
	c.parent = id
}

/*
Получает дочерние ID
*/
func (c Post) GetChildrenIds() []msgId {
	ret := make([]msgId, len(c.children))
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
func (c *Post) AddChildrenIds(ids ...msgId) {
	for _, id := range ids {
		c.children[id] = struct{}{}
	}
}

/*
Простейший конструктор начального поста без родителя
*/
func NewInitPost(id msgId) *Post {
	return &Post{id, id, make(map[msgId]struct{})}
}

/*
Конструктор поста с родителем
*/
func NewPost(id msgId, parent msgId) *Post {
	return &Post{id, parent, make(map[msgId]struct{})}
}

/*
Конструктор поста с родителем и дочерними элементами
*/
func NewPostWithChildren(id msgId, parent msgId, children ...msgId) *Post {
	childrenMap := make(map[msgId]struct{}, len(children))
	for _, child := range children {
		childrenMap[child] = struct{}{}
	}
	return &Post{id, parent, childrenMap}
}
