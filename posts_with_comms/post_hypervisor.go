package comments_and_posts

import "../messages"

type PostHypervisor struct {
	initPostList       map[messages.MsgId]CommentPost
	associatedPostList map[messages.MsgId]messages.MsgId
}
