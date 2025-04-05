package main

import (
	"fmt"
	"lesson-03/documentstore"
)

func main() {
	doc1 := documentstore.Document{
		Fields: map[string]documentstore.DocumentField{
			"key": {
				Type:  documentstore.DocumentFieldTypeString,
				Value: "doc1",
			},
			"title": {
				Type:  documentstore.DocumentFieldTypeString,
				Value: "Document 1",
			},
			"views": {
				Type:  documentstore.DocumentFieldTypeNumber,
				Value: 100,
			},
		},
	}

	doc2 := documentstore.Document{
		Fields: map[string]documentstore.DocumentField{
			"key": {
				Type:  documentstore.DocumentFieldTypeString,
				Value: "doc2",
			},
			"description": {
				Type:  documentstore.DocumentFieldTypeString,
				Value: "Second document example",
			},
			"active": {
				Type:  documentstore.DocumentFieldTypeBool,
				Value: true,
			},
		},
	}

	documentstore.Put(doc1)
	documentstore.Put(doc2)

	if doc, ok := documentstore.Get("doc1"); ok {
		fmt.Println("Retrieved doc1:", doc)
	} else {
		fmt.Println("doc1 not found")
	}

	fmt.Println("List documents:")
	for _, d := range documentstore.List() {
		fmt.Println(d)
	}

	if documentstore.Delete("doc1") {
		fmt.Println("doc1 deleted")
	} else {
		fmt.Println("doc1 not found")
	}

	fmt.Println("Documents after deletion:")
	for _, d := range documentstore.List() {
		fmt.Println(d)
	}
}
