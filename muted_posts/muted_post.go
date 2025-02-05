package muted_posts

import (
	msg "github.com/Nikitarsis/posts_and_comments/messages"
)

type MutedPost struct {
	mutedPost map[msg.MsgId]struct{}
}

func (m MutedPost) CanComment(id msg.MsgId) bool {
	_, check := m.mutedPost[id]
	return check
}

func (m *MutedPost) AllowComment(id msg.MsgId) {
	m.mutedPost[id] = struct{}{}
}

func (m *MutedPost) ForbidComment(id msg.MsgId) {
	delete(m.mutedPost, id)
}

func NewMutedPost() MutedPost {
	return MutedPost{
		mutedPost: make(map[msg.MsgId]struct{}),
	}
}
