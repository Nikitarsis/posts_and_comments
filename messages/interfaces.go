package messages

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
