package documentstore

import (
	"errors"
	"fmt"
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
	config    CollectionConfig
	documents map[string]*Document
}

func NewCollection(cfg CollectionConfig) *Collection {
	return &Collection{
		config:    cfg,
		documents: make(map[string]*Document),
	}
}

func (s *Collection) Put(doc Document) error {
	// Потрібно перевірити що документ містить поле `{cfg.PrimaryKey}` типу `string`
	// TODO: Implement
	primaryKey := s.config.PrimaryKey

	keyField, exists := doc.Fields[primaryKey]
	if !exists {
		return fmt.Errorf("%w: missing field %q", ErrDocumentMissingField, primaryKey)
	}

	if keyField.Type != DocumentFieldTypeString {
		return fmt.Errorf("%w: expected type %q", ErrDocumentHasIncorrectTypeField, DocumentFieldTypeString)
	}

	key, ok := keyField.Value.(string)
	if !ok {
		return fmt.Errorf("%w: keyField is not a string", ErrDocumentHasIncorrectTypeField)
	}

	s.documents[key] = &doc
	return nil
}

func (s *Collection) Get(key string) (*Document, error) {
	// TODO: Implement
	doc, ok := s.documents[key]
	if !ok {
		return nil, ErrDocumentNotFound
	}

	return doc, nil
}

func (s *Collection) Delete(key string) error {
	// TODO: Implement
	if _, ok := s.documents[key]; !ok {
		return ErrDocumentNotFound
	}

	delete(s.documents, key)
	return nil
}

func (s *Collection) List() []Document {
	// TODO: Implement
	var docs []Document
	for _, doc := range s.documents {
		docs = append(docs, *doc)
	}
	return docs
}
