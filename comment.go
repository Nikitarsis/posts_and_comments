package comments_and_posts

type Comment struct {
	id       msgId
	parent   msgId
	children map[msgId]struct{}
}

func (c Comment) GetMessageId() msgId {
	return c.id
}
func (c Comment) GetParentId() msgId {
	return c.parent
}

func (c *Comment) SetParentId(id msgId) {
	c.parent = id
}

func (c Comment) GetChildrenIds() []msgId {
	ret := make([]msgId, len(c.children))
	i := 0
	for k := range c.children {
		ret[i] = k
		i++
	}
	return ret
}

func (c *Comment) AddChildrenIds(ids ...msgId) {
	for _, id := range ids {
		c.children[id] = struct{}{}
	}
}

func (c *Comment) RemoveChildrenIds(ids ...msgId) {
	for _, id := range ids {
		delete(c.children, id)
	}
}

func NewComment(id msgId, parent msgId) *Comment {
	return &Comment{id, parent, make(map[msgId]struct{})}
}

func CommentWithChildren(id msgId, parent msgId, children ...msgId) *Comment {
	childrenMap := make(map[msgId]struct{}, len(children))
	for _, child := range children {
		childrenMap[child] = struct{}{}
	}
	return &Comment{id, parent, childrenMap}
}
