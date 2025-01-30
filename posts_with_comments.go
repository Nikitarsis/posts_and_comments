package comments_and_posts

import (
	"fmt"
)

type CommentPost struct {
	post             IPost
	comments         map[msgId]IComment
	allowNewcomments bool
}

func (c *CommentPost) AllowComment() {
	c.allowNewcomments = true
}

func (c *CommentPost) ForbidComment() {
	c.allowNewcomments = false
}

func (c CommentPost) CanComment() bool {
	return c.allowNewcomments
}

func (c CommentPost) getPost() IPost {
	return c.post
}

func (c CommentPost) getComments(ids ...msgId) ([]IComment, error) {
	ret := make([]IComment, len(ids))
	for i, id := range ids {
		comment, check := c.comments[id]
		if !check {
			str_id := fmt.Sprint(id)
			str_postId := fmt.Sprint(c.post.GetMessageId())
			return nil, fmt.Errorf("no comment with id %s under post with id %s", str_id, str_postId)
		}
		ret[i] = comment
	}
	return ret, nil
}

func (c CommentPost) addComments(ids ...msgId) error {
	if !c.CanComment() {
		str_id := fmt.Sprint(c.post.GetMessageId())
		return fmt.Errorf("post doesn't allow to add commentaries%s", str_id)
	}
	c.post.AddCommentsId(ids...)
	for _, id := range ids {
		c.comments[id] = NewComment(id, c.post.GetMessageId())
	}
	return nil
}

func NewCommentPost(id msgId, allowComments bool) *CommentPost {
	post := NewPost(id)
	comments := make(map[msgId]IComment)
	return &CommentPost{post, comments, allowComments}
}

func NewCommentPostWithComments(id msgId, allowComment bool, comments ...IComment) *CommentPost {
	commentMap := make(map[msgId]IComment, len(comments))
	commentIds := make([]msgId, len(comments))
	for i, comment := range comments {
		cid := comment.GetMessageId()
		commentMap[cid] = comment
		commentIds[i] = cid
	}
	post := NewPostWithComment(id, commentIds...)
	return &CommentPost{post, commentMap, allowComment}
}
