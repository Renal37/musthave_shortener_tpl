package user

type User struct {
	ID  string
	New bool
}

func NewUser(id string, new bool) *User {
	return &User{
		ID:  id,
		New: new,
	}
}
