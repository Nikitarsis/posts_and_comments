package comments_and_posts

type Post struct {
	postId      msgId
	commentsIds map[msgId]struct{}
}

func (p Post) GetMessageId() msgId {
	return p.postId
}

func (p *Post) SetMessageId(id msgId) {
	p.postId = id
}

func (p Post) GetCommentsId() []msgId {
	ret := make([]msgId, len(p.commentsIds))
	i := 0
	for k := range p.commentsIds {
		ret[i] = k
		i++
	}
	return ret
}

func (p *Post) AddCommentsId(ids ...msgId) {
	for _, id := range ids {
		p.commentsIds[id] = struct{}{}
	}
}

func (p *Post) RemoveCommentsId(ids ...msgId) {
	for _, id := range ids {
		delete(p.commentsIds, id)
	}
}

func (p Post) GetCommentsNum() uint {
	return uint(len(p.commentsIds))
}

func NewPostWithComment(id msgId, commentIds ...msgId) *Post {
	comMap := make(map[msgId]struct{}, len(commentIds))
	for commentId := range commentIds {
		comMap[msgId(commentId)] = struct{}{}
	}
	return &Post{id, comMap}
}

func NewPost(id msgId) *Post {
	return &Post{id, map[msgId]struct{}{}}
}
