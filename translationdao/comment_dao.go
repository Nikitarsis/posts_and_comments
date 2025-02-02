package translationdao

import (
	"fmt"
	"strconv"
	"strings"
)

type CommentDao struct {
	UserId      uint64   `json:"user_id"`
	PostId      uint64   `json:"post_id"`
	ParentId    uint64   `json:"parent_id"`
	ChildrenIds []uint64 `json:"children_id"`
	Message     *string  `json:"message,omitempty"`
}

func (c CommentDao) MarshalJSON() ([]byte, error) {
	strChildren := make([]string, len(c.ChildrenIds))
	for i, childId := range c.ChildrenIds {
		strChildren[i] = "\"" + strconv.FormatUint(childId, 16) + "\""
	}
	builder := strings.Builder{}
	builder.WriteString("\"comment\":{")
	builder.WriteString(fmt.Sprintf("\"user_id\":\"%s\",", strconv.FormatUint(c.UserId, 16)))
	builder.WriteString(fmt.Sprintf("\"post_id\":\"%s\",", strconv.FormatUint(c.PostId, 16)))
	builder.WriteString(fmt.Sprintf("\"parent_id\":\"%s\",", strconv.FormatUint(c.ParentId, 16)))
	builder.WriteString(fmt.Sprintf("\"children_ids\":[%s]", strings.Join(strChildren, ",")))
	if c.Message != nil {
		builder.WriteString(fmt.Sprintf(",\"message\":\"%s\"", *c.Message))
	}
	builder.WriteString("}")
	return []byte(builder.String()), nil
}
