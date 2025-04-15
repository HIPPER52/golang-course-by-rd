package documentstore

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
)

var (
	ErrCollectionAlreadyExists = errors.New("collection already exists")
	ErrCollectionNotFound      = errors.New("collection not found")
	ErrUnmarshalDump           = errors.New("failed to unmarshal dump into store")
	ErrMarshalStore            = errors.New("failed to marshal store")
	ErrReadDumpFromFile        = errors.New("failed to read dump from file")
	ErrWriteDumpToFile         = errors.New("failed to write dump to file")
)

type Store struct {
	Collections map[string]*Collection `json:"collections"`
}

func NewStore() *Store {
	return &Store{
		Collections: make(map[string]*Collection),
	}
}

func (s *Store) CreateCollection(name string, cfg *CollectionConfig) (*Collection, error) {
	if _, exists := s.Collections[name]; exists {
		return nil, fmt.Errorf("%w", ErrCollectionAlreadyExists)
	}

	collection := NewCollection(*cfg)
	s.Collections[name] = collection

	slog.Info("Collection created", "name", name)
	return collection, nil
}

func (s *Store) GetCollection(name string) (*Collection, error) {
	collection, exists := s.Collections[name]
	if !exists {
		return nil, fmt.Errorf("%w", ErrCollectionNotFound)
	}
	return collection, nil
}

func (s *Store) DeleteCollection(name string) error {
	if _, exists := s.Collections[name]; !exists {
		return fmt.Errorf("%w", ErrCollectionNotFound)
	}

	slog.Info("Collection deleted", "name", name)
	delete(s.Collections, name)
	return nil
}

func NewStoreFromDump(dump []byte) (*Store, error) {
	// Функція повинна створити на проініціалізувати новий `Store`
	// зі всіма колекціями да даними з вхідного дампу.

	// TODO: Implement
	var s Store
	if err := json.Unmarshal(dump, &s); err != nil {
		return nil, fmt.Errorf("%w", ErrUnmarshalDump)
	}
	return &s, nil
}

func (s *Store) Dump() ([]byte, error) {
	// Методи повинен віддати дамп нашого стору який включані дані про колекції та документ

	// TODO: Implement
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("%w", ErrMarshalStore)
	}
	return data, nil
}

// Значення яке повертає метод `store.Dump()` має без помилок оброблятись функцією `NewStoreFromDump`

func NewStoreFromFile(filename string) (*Store, error) {
	// Робить те ж саме що і функція `NewStoreFromDump`, але сам дамп має діставатись з файлу
	// TODO: Implement
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("%w", ErrReadDumpFromFile)
	}

	slog.Info("Store loaded from file", "filename", filename)
	return NewStoreFromDump(data)
}

func (s *Store) DumpToFile(filename string) error {
	// Робить те ж саме що і метод  `Dump`, але записує у файл замість того щоб повертати сам дамп

	// TODO: Implement
	data, err := s.Dump()
	if err != nil {
		return err
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("%w", ErrWriteDumpToFile)
	}

	slog.Info("Store dumped to file", "filename", filename)
	return nil
}
