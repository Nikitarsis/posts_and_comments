package translationdao

import (
	"fmt"
	"strconv"
	"strings"
)

type CommentDao struct {
	UserId      uint64
	PostId      uint64
	ParentId    uint64
	ChildrenIds []uint64
	Message     *string
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
	builder.WriteString(fmt.Sprintf("\"children_ids\":[%s],", strings.Join(strChildren, ",")))
	builder.WriteString(fmt.Sprintf("\"message\":\"%s\"", *c.Message))
	builder.WriteString("}")
	return []byte(builder.String()), nil
}
