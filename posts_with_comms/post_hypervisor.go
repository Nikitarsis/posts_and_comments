package comments_and_posts

import "github.com/Nikitarsis/posts_and_comments/messages"

type PostHypervisor struct {
	initPostList       map[messages.MsgId]CommentPost
	associatedPostList map[messages.MsgId]messages.MsgId
}
