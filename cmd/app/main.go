package main

import (
	"fmt"
	"lesson_09/internal/document_store"
	"log/slog"
	"os"
)

func init() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)
}

func main() {
	collection := documentstore.NewCollection(documentstore.CollectionConfig{
		PrimaryKey: "id",
	})

	docs := []documentstore.Document{
		{Fields: map[string]documentstore.DocumentField{
			"id":   {Type: documentstore.DocumentFieldTypeString, Value: "1"},
			"name": {Type: documentstore.DocumentFieldTypeString, Value: "Alice"},
		}},
		{Fields: map[string]documentstore.DocumentField{
			"id":   {Type: documentstore.DocumentFieldTypeString, Value: "2"},
			"name": {Type: documentstore.DocumentFieldTypeString, Value: "Bob"},
		}},
		{Fields: map[string]documentstore.DocumentField{
			"id":   {Type: documentstore.DocumentFieldTypeString, Value: "3"},
			"name": {Type: documentstore.DocumentFieldTypeString, Value: "John"},
		}},
	}

	for _, doc := range docs {
		if err := collection.Put(doc); err != nil {
			slog.Error("Failed to put document", "error", err)
		}
	}

	if err := collection.CreateIndex("name"); err != nil {
		slog.Error("Failed to create index", "error", err)
	}

	queryParams := documentstore.QueryParams{
		Desc: false,
	}
	resultDocs, err := collection.Query("name", queryParams)
	if err != nil {
		slog.Error("Failed to query documents", "error", err)
		return
	}

	fmt.Println("Documents sorted by name:")
	for _, doc := range resultDocs {
		nameField := doc.Fields["name"]
		fmt.Printf("- Name: %v\n", nameField.Value)
	}
}
