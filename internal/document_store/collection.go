package documentstore

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
)

var (
	ErrDocumentNotFound              = errors.New("document not found")
	ErrDocumentMissingField          = errors.New("document does not have a key field")
	ErrDocumentHasIncorrectTypeField = errors.New("document has a incorrect type field")
	ErrIndexAlreadyExists            = errors.New("index already exists")
	ErrIndexNotFound                 = errors.New("index not found")
)

type CollectionConfig struct {
	PrimaryKey string
}

type Collection struct {
	mx        sync.RWMutex
	indexesMx sync.RWMutex

	config    CollectionConfig     `json:"config"`
	documents map[string]*Document `json:"documents"`
	indexes   map[string]*Index    `json:"-"`
}

func NewCollection(cfg CollectionConfig) *Collection {
	return &Collection{
		config:    cfg,
		documents: make(map[string]*Document),
		indexes:   make(map[string]*Index),
	}
}

type QueryParams struct {
	Desc     bool    // Визначає в якому порядку повертати дані
	MinValue *string // Визначає мінімальне значення поля для фільтрації
	MaxValue *string // Визначає максимальне значення поля для фільтрації
}

func (s *Collection) CreateIndex(fieldName string) error {
	s.indexesMx.Lock()
	defer s.indexesMx.Unlock()

	if _, exists := s.indexes[fieldName]; exists {
		return ErrIndexAlreadyExists
	}

	idx := &Index{}
	s.indexes[fieldName] = idx

	s.mx.RLock()
	defer s.mx.RUnlock()

	for _, doc := range s.documents {
		field, exists := doc.Fields[fieldName]
		if !exists || field.Type != DocumentFieldTypeString {
			continue
		}
		strValue, ok := field.Value.(string)
		if !ok {
			continue
		}
		idx.Insert(strValue, doc)
	}

	return nil
}

func (s *Collection) DeleteIndex(fieldName string) error {
	s.indexesMx.Lock()
	defer s.indexesMx.Unlock()

	if s.indexes == nil {
		return ErrIndexNotFound
	}

	if _, exists := s.indexes[fieldName]; !exists {
		return ErrIndexNotFound
	}

	delete(s.indexes, fieldName)
	return nil
}

func (s *Collection) Query(fieldName string, params QueryParams) ([]Document, error) {
	s.indexesMx.Lock()
	defer s.indexesMx.Unlock()

	if s.indexes == nil {
		return nil, ErrIndexNotFound
	}

	idx, exists := s.indexes[fieldName]
	if !exists {
		return nil, ErrIndexNotFound
	}

	docsPtrs := idx.RangeQuery(params.MinValue, params.MaxValue, params.Desc)

	var docs []Document = make([]Document, 0, len(docsPtrs))
	for _, ptr := range docsPtrs {
		docs = append(docs, *ptr)
	}

	return docs, nil
}

func (s *Collection) Put(doc Document) error {
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

	s.mx.Lock()
	s.documents[key] = &doc
	s.mx.Unlock()
	return nil
}

func (s *Collection) Get(key string) (*Document, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	doc, ok := s.documents[key]
	if !ok {
		return nil, ErrDocumentNotFound
	}

	return doc, nil
}

func (s *Collection) Delete(key string) error {
	s.mx.Lock()
	defer s.mx.Unlock()

	if _, ok := s.documents[key]; !ok {
		return ErrDocumentNotFound
	}

	delete(s.documents, key)
	return nil
}

func (s *Collection) List() []Document {
	s.mx.RLock()
	defer s.mx.RUnlock()

	docs := make([]Document, 0, len(s.documents))
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
