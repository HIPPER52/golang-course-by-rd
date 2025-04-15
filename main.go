package main

import (
	"fmt"
	"lesson-06/documentstore"
	"lesson-06/users"
	"log/slog"
	"os"
)

func init() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)
}

func main() {
	store := documentstore.NewStore()

	usersCfg := &documentstore.CollectionConfig{PrimaryKey: "id"}
	usersColl, err := store.CreateCollection("users", usersCfg)
	if err != nil {
		slog.Error("failed to create users collection", "error", err)
		os.Exit(1)
	}

	userService := users.NewService(usersColl)

	user1 := users.User{ID: "1", Name: "Alex"}
	user2 := users.User{ID: "2", Name: "John"}

	if _, err := userService.CreateUser(user1); err != nil {
		slog.Error("failed to create Alex", "error", err)
	}
	if _, err := userService.CreateUser(user2); err != nil {
		slog.Error("failed to create John", "error", err)
	}

	filename := "store_dump.json"
	if err := store.DumpToFile(filename); err != nil {
		slog.Error("failed to dump store to file", "file", filename, "error", err)
	} else {
		slog.Info("store dumped to file", "file", filename)
	}

	loadedStore, err := documentstore.NewStoreFromFile(filename)
	if err != nil {
		slog.Error("failed to load store from file", "file", filename, "error", err)
		os.Exit(1)
	}
	slog.Info("store loaded from file", "file", filename)

	loadedUsersColl, err := loadedStore.GetCollection("users")
	if err != nil {
		slog.Error("failed to get users collection from loaded store", "error", err)
		os.Exit(1)
	}

	loadedUserService := users.NewService(loadedUsersColl)

	usersList, err := loadedUserService.ListUsers()
	if err != nil {
		slog.Error("failed to list users from loaded store", "error", err)
		os.Exit(1)
	}

	fmt.Println("Users from loaded store:")
	for _, u := range usersList {
		fmt.Printf("User: ID=%s, Name=%s\n", u.ID, u.Name)
	}

	if err := loadedUserService.DeleteUser(user1.ID); err != nil {
		slog.Error("failed to delete user", "userID", user1.ID, "error", err)
	} else {
		slog.Info("user deleted", "userID", user1.ID)
	}

	usersList, err = loadedUserService.ListUsers()
	if err != nil {
		slog.Error("failed to list users after deletion", "error", err)
		os.Exit(1)
	}

	fmt.Println("Users after deletion:")
	for _, u := range usersList {
		fmt.Printf("User: ID=%s, Name=%s\n", u.ID, u.Name)
	}
}
