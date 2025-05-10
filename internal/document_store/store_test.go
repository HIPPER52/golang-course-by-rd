package documentstore_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"lesson_13/internal/document_store"
)

func TestCreateAndGetCollection(t *testing.T) {
	store := documentstore.NewStore()

	cfg := &documentstore.CollectionConfig{PrimaryKey: "id"}
	coll, err := store.CreateCollection("users", cfg)
	require.NoError(t, err)
	require.NotNil(t, coll)

	loaded, err := store.GetCollection("users")
	require.NoError(t, err)
	require.Same(t, coll, loaded)
}

func TestCreateCollection_AlreadyExists(t *testing.T) {
	store := documentstore.NewStore()
	cfg := &documentstore.CollectionConfig{PrimaryKey: "id"}

	_, err := store.CreateCollection("users", cfg)
	require.NoError(t, err)

	_, err = store.CreateCollection("users", cfg)
	require.ErrorIs(t, err, documentstore.ErrCollectionAlreadyExists)
}

func TestDeleteCollection(t *testing.T) {
	store := documentstore.NewStore()
	cfg := &documentstore.CollectionConfig{PrimaryKey: "id"}

	_, _ = store.CreateCollection("users", cfg)

	err := store.DeleteCollection("users")
	require.NoError(t, err)

	_, err = store.GetCollection("users")
	require.ErrorIs(t, err, documentstore.ErrCollectionNotFound)
}

func TestDumpAndLoad(t *testing.T) {
	store := documentstore.NewStore()
	cfg := &documentstore.CollectionConfig{PrimaryKey: "id"}

	users, _ := store.CreateCollection("users", cfg)

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
	require.NoError(t, users.Put(doc))

	dump, err := store.Dump()
	require.NoError(t, err)

	loaded, err := documentstore.NewStoreFromDump(dump)
	require.NoError(t, err)

	loadedUsers, err := loaded.GetCollection("users")
	require.NoError(t, err)

	got, err := loadedUsers.Get("1")
	require.NoError(t, err)
	require.Equal(t, "Alex", got.Fields["name"].Value)
}

func TestDumpToFile_And_LoadFromFile(t *testing.T) {
	store := documentstore.NewStore()
	cfg := &documentstore.CollectionConfig{PrimaryKey: "id"}

	users, _ := store.CreateCollection("users", cfg)

	doc := documentstore.Document{
		Fields: map[string]documentstore.DocumentField{
			"id": {
				Type:  documentstore.DocumentFieldTypeString,
				Value: "1",
			},
			"email": {
				Type:  documentstore.DocumentFieldTypeString,
				Value: "test@example.com",
			},
		},
	}
	require.NoError(t, users.Put(doc))

	file := "test_store_dump.json.gz"
	defer os.Remove(file)

	require.NoError(t, store.DumpToFile(file))

	loaded, err := documentstore.NewStoreFromFile(file)
	require.NoError(t, err)

	loadedUsers, err := loaded.GetCollection("users")
	require.NoError(t, err)

	got, err := loadedUsers.Get("1")
	require.NoError(t, err)
	require.Equal(t, "test@example.com", got.Fields["email"].Value)
}
