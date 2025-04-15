package main

import (
	"fmt"
	"lesson-05/documentstore"
	"lesson-05/users"
	"log"
)

func main() {
	var store = documentstore.NewStore()
	var config = &documentstore.CollectionConfig{
		PrimaryKey: "id",
	}

	collection, err := store.CreateCollection("users", config)
	if err != nil {
		fmt.Printf("Error creating collection: %v\n", err)
	}

	var userService = users.NewService(collection)

	var user1 = users.User{ID: "1", Name: "Alex"}
	var user2 = users.User{ID: "2", Name: "John"}

	createdUser1, err := userService.CreateUser(user1)
	if err != nil {
		fmt.Printf("Error creating user1: %v\n", err)
	} else {
		fmt.Printf("Created user: %+v\n", createdUser1)
	}

	createdUser2, err := userService.CreateUser(user2)
	if err != nil {
		fmt.Printf("Error creating user2: %v\n", err)
	} else {
		fmt.Printf("Created user: %+v\n", createdUser2)
	}

	usersList, err := userService.ListUsers()
	if err != nil {
		fmt.Printf("Error list users: %v\n", err)
	} else {
		fmt.Println("List of users:")
		for _, user := range usersList {
			fmt.Printf("User: %+v\n", user)
		}
	}

	tempUser, err := userService.GetUser("1")
	if err != nil {
		log.Printf("Error getting user with ID '1': %v\n", err)
	} else {
		fmt.Printf("Got user: %+v\n", tempUser)
	}

	if err := userService.DeleteUser(tempUser.ID); err != nil {
		fmt.Printf("Error deleting user with ID %s: %v\n", tempUser.ID, err)
	} else {
		fmt.Printf("Deleted user with ID %s\n", tempUser.ID)
	}

	usersList, err = userService.ListUsers()
	if err != nil {
		fmt.Printf("Error list users after deletion: %v\n", err)
	} else {
		fmt.Println("List of users after deletion:")
		for _, user := range usersList {
			fmt.Printf("User: %+v\n", user)
		}
	}

	if _, err := userService.GetUser(tempUser.ID); err != nil {
		fmt.Printf("Error getting user: %v\n", err)
	}
}
