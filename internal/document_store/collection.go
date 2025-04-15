package documentstore

import (
	"encoding/json"
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
	config    CollectionConfig     `json:"config"`
	documents map[string]*Document `json:"documents"`
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

	slog.Info("Document added/updated", "collection", s.config.PrimaryKey, "key", key)
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

	slog.Info("Document deleted", "collection", s.config.PrimaryKey, "key", key)
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

func (c *Collection) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Config    CollectionConfig     `json:"config"`
		Documents map[string]*Document `json:"documents"`
	}{
		Config:    c.config,
		Documents: c.documents,
	})
}

func (c *Collection) UnmarshalJSON(data []byte) error {
	aux := &struct {
		Config    CollectionConfig     `json:"config"`
		Documents map[string]*Document `json:"documents"`
	}{}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	c.config = aux.Config
	c.documents = aux.Documents
	return nil
}
