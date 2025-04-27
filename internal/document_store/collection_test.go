package documentstore

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPutAndGet(t *testing.T) {
	collection := NewCollection(CollectionConfig{
		PrimaryKey: "id",
	})

	doc := Document{
		Fields: map[string]DocumentField{
			"id": {
				Type:  DocumentFieldTypeString,
				Value: "1",
			},
			"name": {
				Type:  DocumentFieldTypeString,
				Value: "Alex",
			},
		},
	}

	err := collection.Put(doc)
	require.NoError(t, err)

	got, err := collection.Get("1")
	require.NoError(t, err)
	require.Equal(t, "Alex", got.Fields["name"].Value)
}

func TestPut_InvalidPrimaryKey(t *testing.T) {
	collection := NewCollection(CollectionConfig{
		PrimaryKey: "id",
	})

	docMissing := Document{
		Fields: map[string]DocumentField{
			"name": {
				Type:  DocumentFieldTypeString,
				Value: "Alex",
			},
		},
	}

	err := collection.Put(docMissing)
	require.ErrorIs(t, err, ErrDocumentMissingField)

	docBadType := Document{
		Fields: map[string]DocumentField{
			"id": {
				Type:  DocumentFieldTypeNumber,
				Value: 123,
			},
		},
	}

	err = collection.Put(docBadType)
	require.ErrorIs(t, err, ErrDocumentHasIncorrectTypeField)
}

func TestDelete(t *testing.T) {
	collection := NewCollection(CollectionConfig{
		PrimaryKey: "id",
	})

	doc := Document{
		Fields: map[string]DocumentField{
			"id": {
				Type:  DocumentFieldTypeString,
				Value: "2",
			},
		},
	}

	require.NoError(t, collection.Put(doc))

	err := collection.Delete("2")
	require.NoError(t, err)

	_, err = collection.Get("2")
	require.ErrorIs(t, err, ErrDocumentNotFound)
}

func TestList(t *testing.T) {
	collection := NewCollection(CollectionConfig{
		PrimaryKey: "id",
	})

	docs := []Document{
		{
			Fields: map[string]DocumentField{
				"id": {
					Type:  DocumentFieldTypeString,
					Value: "a1",
				},
			},
		},
		{
			Fields: map[string]DocumentField{
				"id": {
					Type:  DocumentFieldTypeString,
					Value: "b2",
				},
			},
		},
	}

	for _, d := range docs {
		require.NoError(t, collection.Put(d))
	}

	list := collection.List()
	require.Len(t, list, 2)
}

func TestCreateIndex(t *testing.T) {
	coll := NewCollection(CollectionConfig{PrimaryKey: "id"})

	coll.Put(Document{Fields: map[string]DocumentField{
		"id":   {Type: DocumentFieldTypeString, Value: "1"},
		"name": {Type: DocumentFieldTypeString, Value: "Alice"},
	}})
	coll.Put(Document{Fields: map[string]DocumentField{
		"id":   {Type: DocumentFieldTypeString, Value: "2"},
		"name": {Type: DocumentFieldTypeString, Value: "Bob"},
	}})

	err := coll.CreateIndex("name")
	if err != nil {
		t.Fatalf("expected no error creating index, got: %v", err)
	}

	if _, exists := coll.indexes["name"]; !exists {
		t.Fatalf("expected index for field 'name' to exist")
	}
}

func TestQueryWithIndex(t *testing.T) {
	coll := NewCollection(CollectionConfig{PrimaryKey: "id"})

	coll.Put(Document{Fields: map[string]DocumentField{
		"id":   {Type: DocumentFieldTypeString, Value: "1"},
		"name": {Type: DocumentFieldTypeString, Value: "Alice"},
	}})
	coll.Put(Document{Fields: map[string]DocumentField{
		"id":   {Type: DocumentFieldTypeString, Value: "2"},
		"name": {Type: DocumentFieldTypeString, Value: "Bob"},
	}})

	err := coll.CreateIndex("name")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	params := QueryParams{}
	docs, err := coll.Query("name", params)
	if err != nil {
		t.Fatalf("unexpected error querying index: %v", err)
	}

	if len(docs) != 2 {
		t.Fatalf("expected 2 documents, got %d", len(docs))
	}
}

func TestDeleteIndex(t *testing.T) {
	coll := NewCollection(CollectionConfig{PrimaryKey: "id"})

	err := coll.CreateIndex("name")
	if err != nil {
		t.Fatalf("unexpected error creating index: %v", err)
	}

	err = coll.DeleteIndex("name")
	if err != nil {
		t.Fatalf("unexpected error deleting index: %v", err)
	}

	if _, exists := coll.indexes["name"]; exists {
		t.Fatalf("expected index 'name' to be deleted")
	}
}
