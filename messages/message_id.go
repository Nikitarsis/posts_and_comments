package messages

type MsgId struct {
	uint64
}

func (m MsgId) GetId() uint64 {
	return m.uint64
}

func (m MsgId) IsLess(c IMsgId) bool {
	return m.GetId() < c.GetId()
}

func GetMessageId(msg uint64) MsgId {
	return MsgId{msg}
}
