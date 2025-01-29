package comments_and_posts

type Post struct {
	postId      msgId
	commentsIds []msgId
}

func (p Post) GetMessageId() msgId {
	return p.postId
}

func (p *Post) SetMessageId(id msgId) {
	m.messageId = id
}
