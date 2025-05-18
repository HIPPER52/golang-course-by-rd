package documentstore

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"sync"
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
	mx          sync.RWMutex
	Collections map[string]*Collection `json:"collections"`
}

func NewStore() *Store {
	return &Store{
		Collections: make(map[string]*Collection),
	}
}

func (s *Store) CreateCollection(name string, cfg *CollectionConfig) (*Collection, error) {
	s.mx.Lock()
	defer s.mx.Unlock()

	if _, exists := s.Collections[name]; exists {
		return nil, fmt.Errorf("%w", ErrCollectionAlreadyExists)
	}

	collection := NewCollection(*cfg)
	s.Collections[name] = collection

	slog.Info("Collection created", "name", name)
	return collection, nil
}

func (s *Store) GetCollection(name string) (*Collection, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	collection, exists := s.Collections[name]
	if !exists {
		return nil, fmt.Errorf("%w", ErrCollectionNotFound)
	}
	return collection, nil
}

func (s *Store) DeleteCollection(name string) error {
	s.mx.Lock()
	defer s.mx.Unlock()

	if _, exists := s.Collections[name]; !exists {
		return fmt.Errorf("%w", ErrCollectionNotFound)
	}

	slog.Info("Collection deleted", "name", name)
	delete(s.Collections, name)
	return nil
}

func (s *Store) ListCollections() []string {
	s.mx.RLock()
	defer s.mx.RUnlock()

	var names []string
	for name := range s.Collections {
		names = append(names, name)
	}
	return names
}

func (s *Store) ListCollectionNames() []string {
	s.mx.RLock()
	defer s.mx.RUnlock()

	names := make([]string, 0, len(s.Collections))
	for name := range s.Collections {
		names = append(names, name)
	}
	return names
}

func NewStoreFromDump(dump []byte) (*Store, error) {
	var s Store
	if err := json.Unmarshal(dump, &s); err != nil {
		return nil, fmt.Errorf("%w", ErrUnmarshalDump)
	}
	return &s, nil
}

func (s *Store) Dump() ([]byte, error) {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("%w", ErrMarshalStore)
	}
	return data, nil
}

func NewStoreFromFile(filename string) (*Store, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("%w", ErrReadDumpFromFile)
	}

	slog.Info("Store loaded from file", "filename", filename)
	return NewStoreFromDump(data)
}

func (s *Store) DumpToFile(filename string) error {
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
