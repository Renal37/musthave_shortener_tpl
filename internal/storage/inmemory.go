package storage

// Storage представляет собой хранилище URL-адресов с ключами и значениями.
type Storage struct {
	URLs map[string]string // Карта для хранения URL-адресов, где ключом является строка, а значением — соответствующий URL.
}

// NewStorage создаёт и возвращает новый экземпляр хранилища с инициализированной картой URL-адресов.
func NewStorage() *Storage {
	return &Storage{
		URLs: make(map[string]string),
	}
}

// Set добавляет значение value в хранилище по заданному ключу key.
func (s *Storage) Set(key string, value string) {
	s.URLs[key] = value
}

// Get возвращает значение, связанное с заданным ключом key, и флаг наличия этого ключа в хранилище.
func (s *Storage) Get(key string) (string, bool) {
	value, exists := s.URLs[key]
	return value, exists
}
