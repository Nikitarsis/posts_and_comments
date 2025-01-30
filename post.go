package comments_and_posts

import "cmp"

type Post struct {
	id       msgId
	parent   msgId
	children map[msgId]struct{}
}

func (c Post) GetMessageId() msgId {
	return c.id
}
func (c Post) GetParentId() (msgId, bool) {
	return c.parent, cmp.Compare(c.id, c.parent) == 0
}

func (c *Post) SetParentId(id msgId) {
	c.parent = id
}

func (c Post) GetChildrenIds() []msgId {
	ret := make([]msgId, len(c.children))
	i := 0
	for k := range c.children {
		ret[i] = k
		i++
	}
	return ret
}

func (c *Post) AddChildrenIds(ids ...msgId) {
	for _, id := range ids {
		c.children[id] = struct{}{}
	}
}

func (c *Post) RemoveChildrenIds(ids ...msgId) {
	for _, id := range ids {
		delete(c.children, id)
	}
}

func NewInitPost(id msgId) *Post {
	return &Post{id, id, make(map[msgId]struct{})}
}

func NewPost(id msgId, parent msgId) *Post {
	return &Post{id, parent, make(map[msgId]struct{})}
}

func NewPostWithChildren(id msgId, parent msgId, children ...msgId) *Post {
	childrenMap := make(map[msgId]struct{}, len(children))
	for _, child := range children {
		childrenMap[child] = struct{}{}
	}
	return &Post{id, parent, childrenMap}
}
