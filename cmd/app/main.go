package main

import (
	"fmt"
	"lesson_11/internal/document_store"
	"log/slog"
	"math/rand"
	"os"
	"strconv"
	"sync"
)

func init() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)
}

func main() {
	store := documentstore.NewStore()
	usersColl, err := store.CreateCollection("users", &documentstore.CollectionConfig{PrimaryKey: "id"})
	if err != nil {
		slog.Error("failed to create collection", "error", err)
		os.Exit(1)
	}

	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			userID := strconv.Itoa(i)
			doc := documentstore.Document{
				Fields: map[string]documentstore.DocumentField{
					"id": {
						Type:  documentstore.DocumentFieldTypeString,
						Value: userID,
					},
					"name": {
						Type:  documentstore.DocumentFieldTypeString,
						Value: "User" + userID,
					},
				},
			}

			action := rand.Intn(3)
			switch action {
			case 0:
				if err := usersColl.Put(doc); err != nil {
					slog.Warn("put failed", "id", userID, "error", err)
				}
			case 1:
				if _, err := usersColl.Get(userID); err != nil {
					slog.Warn("get failed", "id", userID, "error", err)
				}
			case 2:
				if err := usersColl.Delete(userID); err != nil {
					slog.Warn("delete failed", "id", userID, "error", err)
				}
			}
		}(i)
	}

	wg.Wait()
	fmt.Println("Finished")
}
