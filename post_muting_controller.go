package comments_and_posts

/*
Блокирует ли запрет писать под постом возможность отвечать на ранее оставленные комментарии
Возможно, имеет смысл разрешать так делать.
*/
const ALLOW_SUBCOMMENTS_BLOCK bool = false

type MutingManager struct {
	mutedMsg map[msgId]struct{}
}

func (m MutingManager) CanComment(id msgId) bool {
	_, ret := m.mutedMsg[id]
	return ret
}

func (m *MutingManager) AllowComment(id msgId) {
	m.mutedMsg[id] = struct{}{}
}

func (m *MutingManager) ForbidComment(id msgId) {
	delete(m.mutedMsg, id)
}

func NewMutingManager() *MutingManager {
	return &MutingManager{make(map[msgId]struct{})}
}

func NewMutingManagerBuffered(ids ...msgId) *MutingManager {
	ret := make(map[msgId]struct{}, len(ids))
	for _, id := range ids {
		ret[id] = struct{}{}
	}
	return &MutingManager{ret}
}
