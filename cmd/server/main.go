package main

import (
	"fmt"
	"lesson_14/internal/api"
	"lesson_14/internal/mongodb"
	"net/http"
	"os"
)

func main() {
	fmt.Println("Start launching...")

	if err := mongodb.Init(); err != nil {
		panic(fmt.Errorf("failed to init MongoDB: %v", err))
	}

	dbName := os.Getenv("MONGODB_DB")
	if dbName == "" {
		dbName = "documentstore"
	}

	store := mongodb.NewStore(mongodb.GetClient().Database(dbName))

	handler := api.NewHandler(store)

	http.HandleFunc("/put_document", handler.PutDocument)
	http.HandleFunc("/get_document", handler.GetDocument)
	http.HandleFunc("/list_documents", handler.ListDocuments)
	http.HandleFunc("/delete_document", handler.DeleteDocument)

	http.HandleFunc("/create_collection", handler.CreateCollection)
	http.HandleFunc("/delete_collection", handler.DeleteCollection)
	http.HandleFunc("/list_collections", handler.ListCollections)

	http.HandleFunc("/create_index", handler.CreateIndex)
	http.HandleFunc("/delete_index", handler.DeleteIndex)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(fmt.Errorf("server listening failed: %v", err))
	}
}
