package documentstore

import (
	"strconv"
	"testing"
)

// 358977 ns/op
func BenchmarkQueryWithoutIndex(b *testing.B) {
	coll := NewCollection(CollectionConfig{PrimaryKey: "id"})

	for i := 0; i < 10000; i++ {
		coll.Put(Document{Fields: map[string]DocumentField{
			"id":   {Type: DocumentFieldTypeString, Value: strconv.Itoa(i)},
			"name": {Type: DocumentFieldTypeString, Value: "User" + strconv.Itoa(i)},
		}})
	}

	min := "User2000"
	max := "User3000"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var result []Document
		for _, doc := range coll.List() {
			nameField := doc.Fields["name"]
			if nameStr, ok := nameField.Value.(string); ok {
				if nameStr >= min && nameStr <= max {
					result = append(result, doc)
				}
			}
		}
	}
}

// 38620 ns/op
func BenchmarkQueryWithIndex(b *testing.B) {
	coll := NewCollection(CollectionConfig{PrimaryKey: "id"})

	for i := 0; i < 10000; i++ {
		coll.Put(Document{Fields: map[string]DocumentField{
			"id":   {Type: DocumentFieldTypeString, Value: strconv.Itoa(i)},
			"name": {Type: DocumentFieldTypeString, Value: "User" + strconv.Itoa(i)},
		}})
	}

	coll.CreateIndex("name")

	min := "User2000"
	max := "User3000"

	params := QueryParams{
		MinValue: &min,
		MaxValue: &max,
		Desc:     false,
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = coll.Query("name", params)
	}
}
