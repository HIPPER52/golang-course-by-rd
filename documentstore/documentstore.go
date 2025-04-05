package documentstore

type DocumentFieldType string

const (
	DocumentFieldTypeString DocumentFieldType = "string"
	DocumentFieldTypeNumber DocumentFieldType = "number"
	DocumentFieldTypeBool   DocumentFieldType = "bool"
	DocumentFieldTypeArray  DocumentFieldType = "array"
	DocumentFieldTypeObject DocumentFieldType = "object"
)

type DocumentField struct {
	Type  DocumentFieldType
	Value interface{}
}

type Document struct {
	Fields map[string]DocumentField
}

var documents = make(map[string]Document)

func Put(doc Document) {
	// 1. Перевірити що документ містить в мапі поле `key` типу `string`
	// TODO: Implement
	keyField, exists := doc.Fields["key"]
	if !exists || keyField.Type != DocumentFieldTypeString {
		return
	}

	key, ok := keyField.Value.(string)
	if !ok {
		return
	}

	documents[key] = doc
}

func Get(key string) (*Document, bool) {
	// Потрібно повернути документ по ключу
	// Якщо документ знайдено, повертаємо `true` та поінтер на документ
	// Інакше повертаємо `false` та `nil`
	// TODO: Implement
	doc, ok := documents[key]
	return &doc, ok
}

func Delete(key string) bool {
	// Видаляємо документа по ключу.
	// Повертаємо `true` якщо ми знайшли і видалили документі
	// Повертаємо `false` якщо документ не знайдено
	// TODO: Implement
	if _, ok := documents[key]; ok {
		delete(documents, key)
		return true
	}
	return false
}

func List() []Document {
	// Повертаємо список усіх документів
	// TODO: Implement
	var docs []Document
	for _, doc := range documents {
		docs = append(docs, doc)
	}
	return docs
}
