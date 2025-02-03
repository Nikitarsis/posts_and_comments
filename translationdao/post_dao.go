package translationdao

type PostDao struct {
	PostId  uint64  `json:"post_id"`
	UserId  uint64  `json:"user_id"`
	Message *string `json:"message,omitempty"`
}

func GetEmptyPostDao() PostDao {
	return PostDao{
		PostId:  0,
		UserId:  0,
		Message: nil,
	}
}
