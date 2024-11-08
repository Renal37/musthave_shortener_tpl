package user

// User представляет информацию о пользователе, включая его уникальный идентификатор и статус новизны.
type User struct {
	ID  string // Уникальный идентификатор пользователя.
	New bool   // Флаг, указывающий на то, является ли пользователь новым.
}

// NewUser создаёт и возвращает экземпляр User с заданными ID и статусом новизны.
func NewUser(id string, new bool) *User {
	return &User{
		ID:  id,
		New: new,
	}
}
