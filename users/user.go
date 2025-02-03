package users

type UserId struct {
	uint64
}

func (u UserId) GetId() uint64 {
	return u.uint64
}

type IUser interface {
	GetId() uint64
}

func GetNullUsr() IUser {
	return UserId{uint64: 0}
}

func GetUser(uid uint64) UserId {
	return UserId{uid}
}
