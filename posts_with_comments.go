package comments_and_posts

import (
	"fmt"
)

/*
Реализация поста с комментарием
*/
type CommentPost struct {
	post     IPost
	comments map[msgId]IPost
	//TODO: Подумать над выделенеием в отдельную область
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

/*
Возвращает сам пост
*/
func (c CommentPost) getPost() IPost {
	return c.post
}

/*
Получает комментарии по ID, возвращает nil и ошибку, если ID комментария нет
*/
func (c CommentPost) getComments(ids ...msgId) ([]IPost, error) {
	ret := make([]IPost, len(ids))
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

/*
Добавляет комментарии, возвращает ошибку, если комментировать нельзя
*/
func (c CommentPost) addComments(ids ...msgId) error {
	if !c.CanComment() {
		str_id := fmt.Sprint(c.post.GetMessageId())
		return fmt.Errorf("post doesn't allow to add commentaries %s", str_id)
	}
	c.post.AddChildrenIds(ids...)
	for _, id := range ids {
		c.comments[id] = NewPost(id, c.post.GetMessageId())
	}
	return nil
}

/*
Простой конструктор
*/
func NewCommentPost(id msgId, allowComments bool) *CommentPost {
	post := NewInitPost(id)
	comments := make(map[msgId]IPost)
	return &CommentPost{post, comments, allowComments}
}

/*
Конструктор, принимающий пост и все дочерние элементы.
Возвращает nil и ошибку, если что-либо не сходится
Использовать только при загрузке из базы данных: прожорливый.
*/
func NewCommentPostWithComments(id msgId, allowComment bool, comments ...IPost) *CommentPost {
	commentMap := make(map[msgId]IPost, len(comments))
	commentIds := make([]msgId, len(comments))
	for i, comment := range comments {
		cid := comment.GetMessageId()
		commentMap[cid] = comment
		commentIds[i] = cid
	}
	post := NewPostWithChildren(id, id, commentIds...)
	return &CommentPost{post, commentMap, allowComment}
}
