package messages

import "time"

var hasTimer bool
var lastTime uint64
var counter uint16

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

func GetNewMessageId() MsgId {
	sec := uint64(time.Now().UnixMilli())
	ret := sec << 22
	var max uint64 = 1 << 63
	ret = ret &^ max
	ret = ret | uint64(counter)
	if !hasTimer {
		lastTime = sec
		counter = 0
		return MsgId{ret}
	}
	if lastTime != sec {
		lastTime = sec
		counter = 0
		return MsgId{ret}
	}
	counter++
	return MsgId{ret}
}
