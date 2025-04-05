package main

import (
	"fmt"
	"lesson-04/documentstore"
)

func main() {
	store := documentstore.NewStore()

	cfg := &documentstore.CollectionConfig{
		PrimaryKey: "key",
	}

	created, collection := store.CreateCollection("testCollection", cfg)
	if !created {
		fmt.Println("Collection already exists.")
		return
	}

	doc1 := documentstore.Document{
		Fields: map[string]documentstore.DocumentField{
			"key": {
				Type:  documentstore.DocumentFieldTypeString,
				Value: "doc1",
			},
			"title": {
				Type:  documentstore.DocumentFieldTypeString,
				Value: "First Document",
			},
			"views": {
				Type:  documentstore.DocumentFieldTypeNumber,
				Value: 100,
			},
		},
	}

	collection.Put(doc1)

	doc2 := documentstore.Document{
		Fields: map[string]documentstore.DocumentField{
			"key": {
				Type:  documentstore.DocumentFieldTypeString,
				Value: "doc2",
			},
			"description": {
				Type:  documentstore.DocumentFieldTypeString,
				Value: "Second Document",
			},
			"active": {
				Type:  documentstore.DocumentFieldTypeBool,
				Value: true,
			},
		},
	}

	collection.Put(doc2)

	if d, ok := collection.Get("doc1"); ok {
		fmt.Println("Retrieved doc1: ", d)
	} else {
		fmt.Println("doc1 not found")
	}

	fmt.Println("All documents in collection:")
	for _, d := range collection.List() {
		fmt.Println(d)
	}

	if collection.Delete("doc1") {
		fmt.Println("doc1 deleted")
	} else {
		fmt.Println("doc1 not found")
	}

	fmt.Println("Documents after deletion:")
	for _, d := range collection.List() {
		fmt.Println(d)
	}

	if c, ok := store.GetCollection("testCollection"); ok {
		fmt.Println("Got collection from store:", c)
	} else {
		fmt.Println("Collection not found")
	}

	if store.DeleteCollection("testCollection") {
		fmt.Println("Collection deleted from store")
	} else {
		fmt.Println("Failed to delete collection from store")
	}
}
