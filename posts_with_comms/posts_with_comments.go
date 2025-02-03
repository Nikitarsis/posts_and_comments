package comments_and_posts

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"sync"

	"github.com/Nikitarsis/posts_and_comments/messages"
)

/*
Реализация поста с комментарием
*/
type CommentPost struct {
	post         IPost
	comments     map[messages.MsgId]IPost
	commentPages []messages.MsgId
	mutex        *sync.Mutex
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
func (c CommentPost) getComments(ids ...messages.MsgId) ([]IPost, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
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

func (c CommentPost) getCommentPage(from int, to int) ([]messages.MsgId, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if from < 0 {
		return nil, errors.New("too little")
	}
	if to >= len(c.commentPages) {
		return nil, errors.New("too big")
	}
	return c.commentPages[from:to], nil
}

/*
Добавляет комментарии непосредственно к посту
*/
func (c *CommentPost) addCommentsToPost(ids ...messages.MsgId) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.post.AddChildrenIds(ids...)
	for _, id := range ids {
		c.comments[id] = NewPost(id, c.post.GetMessageId())
	}
	c.commentPages = append(c.commentPages, ids...)
	sort.Slice(c.commentPages, func(i int, j int) bool {
		return c.commentPages[i].IsLess(c.commentPages[j])
	})
}

/*
Добавляет побочные комментарии, возвращает ошибку, если не находится родительского комментария
*/
func (c *CommentPost) addSubcomments(commentId messages.MsgId, ids ...messages.MsgId) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	comment, check := c.comments[commentId]
	if !check {
		str_id := fmt.Sprint(c.post.GetMessageId())
		str_cid := fmt.Sprint(commentId)
		return fmt.Errorf("comment %s is not belong to post %s", str_cid, str_id)
	}
	comment.AddChildrenIds(ids...)
	c.commentPages = append(c.commentPages, ids...)
	sort.Slice(c.commentPages, func(i int, j int) bool {
		return c.commentPages[i].IsLess(c.commentPages[j])
	})
	return nil
}

/*
Простой конструктор
*/
func NewCommentPost(id messages.MsgId) *CommentPost {
	post := NewInitPost(id)
	comments := make(map[messages.MsgId]IPost)
	var mutex sync.Mutex
	sortedSlice := make([]messages.MsgId, 0)
	return &CommentPost{post, comments, sortedSlice, &mutex}
}

/*
Конструктор, принимающий пост и все дочерние элементы.
Возвращает nil и ошибку, если что-либо не сходится
Использовать только при загрузке из базы данных.
*/
//TODO Если останется время, переделать ещё раз
func NewCommentPostWithComments(post IPost, comments ...IPost) (*CommentPost, error) {
	//Создание списка комментариев
	comMap := make(map[messages.MsgId]IPost, len(comments))
	sortSlice := make([]messages.MsgId, len(comments))
	for i, comment := range comments {
		id := comment.GetMessageId()
		comMap[id] = comment
		sortSlice[i] = id
	}
	sort.Slice(sortSlice, func(i, j int) bool {
		return sortSlice[i].IsLess(sortSlice[j])
	})
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
				err := fmt.Errorf("subcomment %s under comment %s has no place", str_chid, str_cid)
				return nil, err
			}
		}
	}
	var mutex sync.Mutex
	return &CommentPost{post, comMap, sortSlice, &mutex}, nil
}
