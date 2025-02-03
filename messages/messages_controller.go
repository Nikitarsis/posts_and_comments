package messages

type MessagesController struct {
	content map[MsgId]*string
}

func (m MessagesController) GetContent(id MsgId) *string {
	ret, check := m.content[id]
	if !check {
		str := ""
		return &str
	}
	return ret
}

func (m *MessagesController) SetContent(id MsgId, message *string) {
	m.content[id] = message
}

func (m *MessagesController) DeleteContent(id MsgId) {
	delete(m.content, id)
}
