package main

import (
	"context"
	"fmt"
	"lesson_14/internal/api"
	"lesson_14/internal/mongodb"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	fmt.Println("Start launching...")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	client, err := mongodb.NewMongoClient(ctx)
	if err != nil {
		panic(fmt.Errorf("failed to init MongoDB: %v", err))
	}

	dbName := os.Getenv("MONGODB_DB")
	if dbName == "" {
		dbName = "documentstore"
	}

	store := mongodb.NewStore(client.Database(dbName))

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

	srv := &http.Server{Addr: ":" + port}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
		}
	}()

	<-ctx.Done()
	slog.Info("Shutdown initiated")

	if err := client.Disconnect(context.Background()); err != nil {
		slog.Error("Failed to disconnect MongoDB", "error", err)
	}
}
