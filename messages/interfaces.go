package messages

// Обложка для сообщения
type content string

type MsgId struct {
	uint64
}

// Структура имеет ID сообщения
type IHaveMessageId interface {
	//Возвращает свой ID
	GetMessageId() MsgId
}

// Структура содержит содержание(текст, ссылки на прикреплёные файлы etc.)
type IHaveContent interface {
	//Возвращает содержание
	GetContent() content
	//Задаёт содержание
	SetContent(data content)
}

// Сообщение с Id и содержанием
type IMessage interface {
	IHaveMessageId
	IHaveContent
}
