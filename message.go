package comments_and_posts

type POST_TYPE string

const (
	POST    POST_TYPE = "post"
	COMMENT POST_TYPE = "comment"
)

type Message struct {
	messageId      msgId
	contentMessage content
}

func (m Message) GetMessageId() msgId {
	return m.messageId
}

func (m Message) GetContent() content {
	return m.contentMessage
}

func (m *Message) SetContent(data content) {
	m.contentMessage = data
}

func NewMessage(id msgId, data content) *Message {
	return &Message{id, data}
}
