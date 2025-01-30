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

func (p *Post) AddCommentId(ids ...msgId) {
	for _, id := range ids {
		p.commentsIds[id] = struct{}{}
	}
}

func (p *Post) RemoveCommentId(ids ...msgId) {
	for _, id := range ids {
		delete(p.commentsIds, id)
	}
}

func NewPost(id msgId, commentIds ...msgId) *Post {
	comMap := make(map[msgId]struct{}, len(commentIds))
	for commentId := range commentIds {
		comMap[msgId(commentId)] = struct{}{}
	}
	return &Post{id, comMap}
}
