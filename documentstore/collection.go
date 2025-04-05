package documentstore

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

func (s *Collection) Put(doc Document) {
	// Потрібно перевірити що документ містить поле `{cfg.PrimaryKey}` типу `string`
	// TODO: Implement
	primaryKey := s.config.PrimaryKey

	keyField, exists := doc.Fields[primaryKey]
	if !exists || keyField.Type != DocumentFieldTypeString {
		return
	}

	key, ok := keyField.Value.(string)
	if !ok {
		return
	}

	s.documents[key] = &doc
}

func (s *Collection) Get(key string) (*Document, bool) {
	// TODO: Implement
	doc, ok := s.documents[key]
	return doc, ok
}

func (s *Collection) Delete(key string) bool {
	// TODO: Implement
	if _, ok := s.documents[key]; ok {
		delete(s.documents, key)
		return true
	}
	return false
}

func (s *Collection) List() []Document {
	// TODO: Implement
	var docs []Document
	for _, doc := range s.documents {
		docs = append(docs, *doc)
	}
	return docs
}
