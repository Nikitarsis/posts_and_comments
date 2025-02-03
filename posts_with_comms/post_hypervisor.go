package comments_and_posts

import (
	"github.com/Nikitarsis/posts_and_comments/messages"
)

type PostHypervisor struct {
	initPostList       map[messages.MsgId]CommentPost
	associatedPostList map[messages.MsgId]messages.MsgId
}

func (p PostHypervisor) GetPosts() []messages.MsgId {
	ret := make([]messages.MsgId, len(p.initPostList))
	i := 0
	for key := range p.initPostList {
		ret[i] = key
		i++
	}
	return ret
}

func (p PostHypervisor) GetPost(id messages.MsgId) (*CommentPost, bool) {
	ret, check := p.initPostList[id]
	if !check {
		return nil, false
	}
	return &ret, true
}

func (p *PostHypervisor) NewPost(post messages.MsgId) {
	comPost := NewCommentPost(post)
	p.initPostList[post] = *comPost
	p.associatedPostList[post] = post
}

func (p PostHypervisor) HasPost(post messages.MsgId) bool {
	_, check := p.initPostList[post]
	return check
}

func NewPostHypervisor() PostHypervisor {
	return PostHypervisor{
		initPostList:       make(map[messages.MsgId]CommentPost),
		associatedPostList: make(map[messages.MsgId]messages.MsgId),
	}
}
