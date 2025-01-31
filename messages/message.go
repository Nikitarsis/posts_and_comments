package messages

// Тип сообщения
type POST_TYPE string

const (
	POST    POST_TYPE = "post"
	COMMENT POST_TYPE = "comment"
)

/*
Сообщение внутри комментария, поста etc. Его содержательная часть
Содержит уникальный идентифекатор и само содержание
*/
type Message struct {
	messageId      MsgId
	contentMessage content
}

/*
Возвращает ID сообщения
*/
func (m Message) GetMessageId() MsgId {
	return m.messageId
}

/*
Получает содержание сообщения
*/
func (m Message) GetContent() content {
	return m.contentMessage
}

/*
Обновляет содержание сообщения
*/
func (m *Message) SetContent(data content) {
	m.contentMessage = data
}

/*
Простой конструктор
*/
func NewMessage(id MsgId, data content) *Message {
	return &Message{id, data}
}
