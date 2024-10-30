package storage

type Storage struct {
    URLs map[string]string
}
// Функция создания нового экземпляра хранилища
func NewStorage() *Storage {
    return &Storage{
        URLs: make(map[string]string),
    }
}
// Функция установки значения по ключу
func (s *Storage) Set(key string, value string) {
    s.URLs[key] = value
}
// Функция получения значения по ключу
func (s *Storage) Get(key string) (string, bool) {
    value, exists := s.URLs[key]
    return value, exists
}
