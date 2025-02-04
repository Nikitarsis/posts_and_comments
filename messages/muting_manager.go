package messages

import (
	"database/sql"
	"fmt"

	sqlc "github.com/Nikitarsis/posts_and_comments/sql_connection"
)

/*
Блокирует ли запрет писать под постом возможность отвечать на ранее оставленные комментарии
Возможно, имеет смысл разрешать так делать.
*/
const ALLOW_SUBCOMMENTS_BLOCK bool = false

type MutingManager struct {
	mutedMsg  map[MsgId]struct{}
	writeToDb func(string)
}

func (m MutingManager) CanComment(id MsgId) bool {
	_, ret := m.mutedMsg[id]
	return ret
}

func (m *MutingManager) AllowComment(id MsgId) {
	m.mutedMsg[id] = struct{}{}
	m.writeToDb(fmt.Sprintf("INSERT INTO mutedDB (id)(%d)", id.GetId()))
}

func (m *MutingManager) ForbidComment(id MsgId) {
	delete(m.mutedMsg, id)
	m.writeToDb(fmt.Sprintf("DELETE FROM mutedDB WHERE id=%d", id.GetId()))
}

func NewMutingManager() *MutingManager {
	return &MutingManager{
		mutedMsg:  make(map[MsgId]struct{}),
		writeToDb: func(s string) {},
	}
}

func NewMutingManagerBuffered(ids ...MsgId) *MutingManager {
	ret := make(map[MsgId]struct{}, len(ids))
	for _, id := range ids {
		ret[id] = struct{}{}
	}
	return &MutingManager{
		mutedMsg:  ret,
		writeToDb: func(s string) {},
	}
}

func NewMutingManagerSQL(connection sqlc.ConnectionSQL) *MutingManager {
	exist := connection.HasTable("mutedDB")
	if !exist {
		connection.AddTable("mutedDB", "(id int8 PRIMARY KEY)")
		ret := NewMutingManager()
		ret.writeToDb = func(s string) { sqlc.Exec(connection, s) }
		return ret
	}
	ids := sqlc.GetObject(
		connection,
		"SELECT * AS id FROM mutedDB",
		func(r *sql.Rows) []MsgId {
			defer r.Close()
			ret := make([]MsgId, 0)
			for r.Next() {
				var id uint64
				r.Scan(&id)
				ret = append(ret, GetMessageId(id))
			}
			return ret
		},
	)
	ret := NewMutingManagerBuffered(ids...)
	ret.writeToDb = func(s string) { sqlc.Exec(connection, s) }
	return ret
}
