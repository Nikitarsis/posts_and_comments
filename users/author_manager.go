package users

import (
	"github.com/Nikitarsis/posts_and_comments/messages"
)

type AuthorManager struct {
	authors        map[UserId]struct{}
	authorsOfPosts map[messages.MsgId]UserId
}

func (a AuthorManager) GetAuthorOfPost(id messages.MsgId) (IUser, bool) {
	ret, check := a.authorsOfPosts[id]
	return ret, check
}

func (a AuthorManager) CheckAuthor(id UserId) bool {
	_, ret := a.authors[id]
	return ret
}

func NewAuthorManager() AuthorManager {
	return AuthorManager{
		authors:        make(map[UserId]struct{}),
		authorsOfPosts: make(map[messages.MsgId]UserId),
	}
}
