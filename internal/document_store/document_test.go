package documentstore_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"lesson_09/internal/document_store"
)

func TestDocumentFieldTypeAssignment(t *testing.T) {
	doc := documentstore.Document{
		Fields: map[string]documentstore.DocumentField{
			"username": {
				Type:  documentstore.DocumentFieldTypeString,
				Value: "hipper",
			},
			"age": {
				Type:  documentstore.DocumentFieldTypeNumber,
				Value: 25,
			},
			"verified": {
				Type:  documentstore.DocumentFieldTypeBool,
				Value: true,
			},
		},
	}

	require.Equal(t, "hipper", doc.Fields["username"].Value)
	require.Equal(t, 25, doc.Fields["age"].Value)
	require.Equal(t, true, doc.Fields["verified"].Value)
}

func TestMarshalUnmarshalDocument(t *testing.T) {
	type User struct {
		ID    string `json:"id"`
		Email string `json:"email"`
	}

	original := User{
		ID:    "u123",
		Email: "hipper@example.com",
	}

	doc, err := documentstore.MarshalDocument(original)
	require.NoError(t, err)

	var result User
	err = documentstore.UnmarshalDocument(doc, &result)
	require.NoError(t, err)

	require.Equal(t, original.ID, result.ID)
	require.Equal(t, original.Email, result.Email)
}
