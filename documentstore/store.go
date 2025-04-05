package documentstore

type Store struct {
	// ...
	collections map[string]*Collection
}

func NewStore() *Store {
	// TODO: Implement
	return &Store{
		collections: make(map[string]*Collection),
	}
}

func (s *Store) CreateCollection(name string, cfg *CollectionConfig) (bool, *Collection) {
	// Створюємо нову колекцію і повертаємо `true` якщо колекція була створена
	// Якщо ж колекція вже створеня то повертаємо `false` та nil
	// TODO: Implement
	if _, exists := s.collections[name]; exists {
		return false, nil
	}

	collection := NewCollection(*cfg)
	s.collections[name] = collection

	return true, collection
}

func (s *Store) GetCollection(name string) (*Collection, bool) {
	// TODO: Implement
	collection, exists := s.collections[name]
	return collection, exists
}

func (s *Store) DeleteCollection(name string) bool {
	// TODO: Implement
	if _, exists := s.collections[name]; exists {
		delete(s.collections, name)
		return true
	}

	return false
}
