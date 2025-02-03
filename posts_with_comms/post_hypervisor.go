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
