package documentstore_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"lesson_07/internal/document_store"
)

func TestPutAndGet(t *testing.T) {
	collection := documentstore.NewCollection(documentstore.CollectionConfig{
		PrimaryKey: "id",
	})

	doc := documentstore.Document{
		Fields: map[string]documentstore.DocumentField{
			"id": {
				Type:  documentstore.DocumentFieldTypeString,
				Value: "1",
			},
			"name": {
				Type:  documentstore.DocumentFieldTypeString,
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
	collection := documentstore.NewCollection(documentstore.CollectionConfig{
		PrimaryKey: "id",
	})

	docMissing := documentstore.Document{
		Fields: map[string]documentstore.DocumentField{
			"name": {
				Type:  documentstore.DocumentFieldTypeString,
				Value: "Alex",
			},
		},
	}

	err := collection.Put(docMissing)
	require.ErrorIs(t, err, documentstore.ErrDocumentMissingField)

	docBadType := documentstore.Document{
		Fields: map[string]documentstore.DocumentField{
			"id": {
				Type:  documentstore.DocumentFieldTypeNumber,
				Value: 123,
			},
		},
	}

	err = collection.Put(docBadType)
	require.ErrorIs(t, err, documentstore.ErrDocumentHasIncorrectTypeField)
}

func TestDelete(t *testing.T) {
	collection := documentstore.NewCollection(documentstore.CollectionConfig{
		PrimaryKey: "id",
	})

	doc := documentstore.Document{
		Fields: map[string]documentstore.DocumentField{
			"id": {
				Type:  documentstore.DocumentFieldTypeString,
				Value: "2",
			},
		},
	}

	require.NoError(t, collection.Put(doc))

	err := collection.Delete("2")
	require.NoError(t, err)

	_, err = collection.Get("2")
	require.ErrorIs(t, err, documentstore.ErrDocumentNotFound)
}

func TestList(t *testing.T) {
	collection := documentstore.NewCollection(documentstore.CollectionConfig{
		PrimaryKey: "id",
	})

	docs := []documentstore.Document{
		{
			Fields: map[string]documentstore.DocumentField{
				"id": {
					Type:  documentstore.DocumentFieldTypeString,
					Value: "a1",
				},
			},
		},
		{
			Fields: map[string]documentstore.DocumentField{
				"id": {
					Type:  documentstore.DocumentFieldTypeString,
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
