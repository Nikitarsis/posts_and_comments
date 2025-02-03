package comments_and_posts

import (
	"github.com/Nikitarsis/posts_and_comments/messages"
)

type PostHypervisor struct {
	initPostList       map[messages.MsgId]*CommentPost
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
	return ret, true
}

func (p PostHypervisor) GetComment(id messages.MsgId) (IPost, bool) {
	parent := p.associatedPostList[id]
	if parent == id {
		return nil, false
	}
	comPost, check := p.initPostList[parent]
	if !check {
		return nil, false
	}
	ret, err := comPost.GetComments(id)
	if err != nil {
		return nil, false
	}
	return ret[0], true
}

func (p *PostHypervisor) NewPost(post messages.MsgId) {
	comPost := NewCommentPost(post)
	p.initPostList[post] = comPost
	p.associatedPostList[post] = post
}

func (p *PostHypervisor) NewComment(parent messages.MsgId, comment messages.MsgId) {
	post := p.associatedPostList[parent]
	p.initPostList[post].AddCommentsToPost(comment)
	p.associatedPostList[comment] = post
}

func (p *PostHypervisor) DeletePost(post messages.MsgId) {
	_, check := p.initPostList[post]
	if !check {
		return
	}
	delete(p.initPostList, post)
	for key, value := range p.associatedPostList {
		if value == post {
			delete(p.associatedPostList, key)
		}
	}
	delete(p.associatedPostList, post)
}

func (p PostHypervisor) HasPost(post messages.MsgId) bool {
	_, check := p.initPostList[post]
	return check
}

func NewPostHypervisor() PostHypervisor {
	return PostHypervisor{
		initPostList:       make(map[messages.MsgId]*CommentPost),
		associatedPostList: make(map[messages.MsgId]messages.MsgId),
	}
}
