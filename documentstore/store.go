package documentstore

import (
	"errors"
	"fmt"
)

var (
	ErrCollectionAlreadyExists = errors.New("collection already exists")
	ErrCollectionNotFound      = errors.New("collection not found")
)

type Store struct {
	collections map[string]*Collection
}

func NewStore() *Store {
	// TODO: Implement
	return &Store{
		collections: make(map[string]*Collection),
	}
}

func (s *Store) CreateCollection(name string, cfg *CollectionConfig) (*Collection, error) {
	// Створюємо нову колекцію і повертаємо `true` якщо колекція була створена
	// Якщо ж колекція вже створеня то повертаємо `false` та nil
	// TODO: Implement
	if _, exists := s.collections[name]; exists {
		return nil, fmt.Errorf("%w", ErrCollectionAlreadyExists)
	}

	collection := NewCollection(*cfg)
	s.collections[name] = collection

	return collection, nil
}

func (s *Store) GetCollection(name string) (*Collection, error) {
	// TODO: Implement
	collection, exists := s.collections[name]
	if !exists {
		return nil, fmt.Errorf("%w", ErrCollectionNotFound)
	}
	return collection, nil
}

func (s *Store) DeleteCollection(name string) error {
	// TODO: Implement
	if _, exists := s.collections[name]; !exists {
		return fmt.Errorf("%w", ErrCollectionNotFound)
	}

	delete(s.collections, name)
	return nil
}
