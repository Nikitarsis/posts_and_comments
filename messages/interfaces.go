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

type IMsgId interface {
	GetId() uint64
	IsLess(c IMsgId) bool
}

// Структура имеет ID сообщения
type IHaveMessageId interface {
	//Возвращает свой ID
	GetMessageId() MsgId
}

// Структура содержит содержание(текст, ссылки на прикреплёные файлы etc.)
type IHaveContent interface {
	//Возвращает содержание
	GetContent() string
	//Задаёт содержание
	SetContent(data string)
}

// Сообщение с Id и содержанием
type IMessage interface {
	IHaveMessageId
	IHaveContent
}
