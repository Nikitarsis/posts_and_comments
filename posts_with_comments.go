package comments_and_posts

import (
	"fmt"
	"reflect"
)

/*
Реализация поста с комментарием
*/
type CommentPost struct {
	post     IPost
	comments map[msgId]IPost
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
Добавляет комментарии непосредственно к посту
*/
func (c *CommentPost) addCommentsToPost(ids ...msgId) {
	c.post.AddChildrenIds(ids...)
	for _, id := range ids {
		c.comments[id] = NewPost(id, c.post.GetMessageId())
	}
}

/*
Добавляет побочные комментарии, возвращает ошибку, если не находится родительского комментария
*/
func (c *CommentPost) addSubcomments(commentId msgId, ids ...msgId) error {
	comment, check := c.comments[commentId]
	if !check {
		str_id := fmt.Sprint(c.post.GetMessageId())
		str_cid := fmt.Sprint(commentId)
		return fmt.Errorf("comment %s is not belong to post %s", str_cid, str_id)
	}
	comment.AddChildrenIds(ids...)
	return nil
}

/*
Простой конструктор
*/
func NewCommentPost(id msgId) *CommentPost {
	post := NewInitPost(id)
	comments := make(map[msgId]IPost)
	return &CommentPost{post, comments}
}

/*
Конструктор, принимающий пост и все дочерние элементы.
Возвращает nil и ошибку, если что-либо не сходится
Использовать только при загрузке из базы данных.
*/
//TODO Если останется время, переделать ещё раз
func NewCommentPostWithComments(post IPost, comments ...IPost) (*CommentPost, error) {
	//Создание списка комментариев
	comMap := make(map[msgId]IPost, len(comments))
	for _, comment := range comments {
		id := comment.GetMessageId()
		comMap[id] = comment
	}
	//Проверка поста на отсутствия родителей
	if _, init := post.GetParentId(); !init {
		str_id := fmt.Sprint(post.GetMessageId())
		err := fmt.Errorf("post with id %s is not init", str_id)
		return nil, err
	}
	//Проверка поста на то, имеются ли его дочерние элементы в comments
	for _, childId := range post.GetChildrenIds() {
		if _, check := comMap[childId]; !check {
			str_pid := fmt.Sprint(post.GetMessageId())
			str_cid := fmt.Sprint(post.GetMessageId())
			return nil, fmt.Errorf("comment %s doesn't belong to post %s", str_cid, str_pid)
		}
	}
	//Проверка всех комментариев на наличие в общем списке родителей и дочерних элементов
	for _, comment := range comments {
		parentId, _ := comment.GetParentId()
		commentId := comment.GetMessageId()
		//Если родителя нет в списке комментариев, проверяется основной пост, и в случае неудачи, возвращается ошибка.
		if _, check := comMap[parentId]; !check {
			if !reflect.DeepEqual(parentId, commentId) {
				str_pid := fmt.Sprint(parentId)
				str_cid := fmt.Sprint(commentId)
				err := fmt.Errorf("comment %s has no place under post %s", str_cid, str_pid)
				return nil, err
			}
		}
		childrenIds := comment.GetChildrenIds()
		//Проверка всех дочерних элементов на предмет наличия в мапе, если ID в мапе нет, возвращается ошибка
		for _, childId := range childrenIds {
			if _, check := comMap[childId]; !check {
				str_chid := fmt.Sprint(childId)
				str_cid := fmt.Sprint(commentId)
				err := fmt.Errorf("Subcomment %s under comment %s has no place", str_chid, str_cid)
				return nil, err
			}
		}
	}
	return &CommentPost{post, comMap}, nil
}
