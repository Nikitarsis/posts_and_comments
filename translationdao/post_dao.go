package translationdao

import (
	"fmt"
	"strconv"
	"strings"
)

type PostDao struct {
	PostId  uint64
	UserId  uint64
	Message *string
}

func (p PostDao) MarshalJSON() ([]byte, error) {
	builder := strings.Builder{}
	builder.WriteString("\"post\":{")
	builder.WriteString(fmt.Sprintf("\"user_id\":\"%s\",", strconv.FormatUint(p.UserId, 16)))
	builder.WriteString(fmt.Sprintf("\"post_id\":\"%s\",", strconv.FormatUint(p.PostId, 16)))
	builder.WriteString(fmt.Sprintf("\"message\":\"%s\"", *p.Message))
	builder.WriteString("}")
	return []byte(builder.String()), nil
}
