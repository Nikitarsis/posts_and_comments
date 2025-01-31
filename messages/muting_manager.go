package messages

/*
Блокирует ли запрет писать под постом возможность отвечать на ранее оставленные комментарии
Возможно, имеет смысл разрешать так делать.
*/
const ALLOW_SUBCOMMENTS_BLOCK bool = false

type MutingManager struct {
	mutedMsg map[MsgId]struct{}
}

func (m MutingManager) CanComment(id MsgId) bool {
	_, ret := m.mutedMsg[id]
	return ret
}

func (m *MutingManager) AllowComment(id MsgId) {
	m.mutedMsg[id] = struct{}{}
}

func (m *MutingManager) ForbidComment(id MsgId) {
	delete(m.mutedMsg, id)
}

func NewMutingManager() *MutingManager {
	return &MutingManager{make(map[MsgId]struct{})}
}

func NewMutingManagerBuffered(ids ...MsgId) *MutingManager {
	ret := make(map[MsgId]struct{}, len(ids))
	for _, id := range ids {
		ret[id] = struct{}{}
	}
	return &MutingManager{ret}
}
