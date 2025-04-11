package documentstore

import (
	"errors"
	"fmt"
	"log/slog"
)

var (
	ErrDocumentNotFound              = errors.New("document not found")
	ErrDocumentMissingField          = errors.New("document does not have a key field")
	ErrDocumentHasIncorrectTypeField = errors.New("document has a incorrect type field")
)

type CollectionConfig struct {
	PrimaryKey string
}

type Collection struct {
	Config    CollectionConfig
	Documents map[string]*Document
}

func NewCollection(cfg CollectionConfig) *Collection {
	return &Collection{
		Config:    cfg,
		Documents: make(map[string]*Document),
	}
}

func (s *Collection) Put(doc Document) error {
	// Потрібно перевірити що документ містить поле `{cfg.PrimaryKey}` типу `string`
	// TODO: Implement
	primaryKey := s.Config.PrimaryKey

	keyField, exists := doc.Fields[primaryKey]
	if !exists {
		return fmt.Errorf("%w: missing field %q", ErrDocumentMissingField, primaryKey)
	}

	if keyField.Type != DocumentFieldTypeString {
		return fmt.Errorf("%w: expected type %q", ErrDocumentHasIncorrectTypeField, DocumentFieldTypeString)
	}

	key, ok := keyField.Value.(string)
	if !ok {
		return fmt.Errorf("%w: value is not a string", ErrDocumentHasIncorrectTypeField)
	}

	slog.Info("Document added/updated", "collection", s.Config.PrimaryKey, "key", key)
	s.Documents[key] = &doc
	return nil
}

func (s *Collection) Get(key string) (*Document, error) {
	// TODO: Implement
	doc, ok := s.Documents[key]
	if !ok {
		return nil, fmt.Errorf("%w", ErrDocumentNotFound)
	}

	return doc, nil
}

func (s *Collection) Delete(key string) error {
	// TODO: Implement
	if _, ok := s.Documents[key]; !ok {
		return fmt.Errorf("%w", ErrDocumentNotFound)
	}

	slog.Info("Document deleted", "collection", s.Config.PrimaryKey, "key", key)
	delete(s.Documents, key)
	return nil
}

func (s *Collection) List() []Document {
	// TODO: Implement
	var docs []Document
	for _, doc := range s.Documents {
		docs = append(docs, *doc)
	}
	return docs
}
