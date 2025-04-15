package users

import (
	"errors"
	"fmt"
	"lesson-06/documentstore"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Service struct {
	coll *documentstore.Collection
}

func NewService(coll *documentstore.Collection) *Service {
	return &Service{coll: coll}
}

func (s *Service) CreateUser(user User) (*User, error) {
	var document, err = documentstore.MarshalDocument(user)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal user: %w", err)
	}

	if err := s.coll.Put(*document); err != nil {
		return nil, fmt.Errorf("failed to put user document: %w", err)
	}

	return &user, nil
}

func (s *Service) ListUsers() ([]User, error) {
	var documents = s.coll.List()
	var users = make([]User, 0, len(documents))
	for _, doc := range documents {
		var user User
		if err := documentstore.UnmarshalDocument(&doc, &user); err != nil {
			return nil, fmt.Errorf("failed to unmarshal document: %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}

func (s *Service) GetUser(userID string) (*User, error) {
	var document, err = s.coll.Get(userID)
	if err != nil {
		return nil, fmt.Errorf("%w: user with id %s", ErrUserNotFound, userID)
	}

	var user User
	if err := documentstore.UnmarshalDocument(document, &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal document into user: %w", err)
	}
	return &user, nil
}

func (s *Service) DeleteUser(userID string) error {
	if err := s.coll.Delete(userID); err != nil {
		return fmt.Errorf("failed to delete user with id %s: %w", userID, err)
	}

	return nil
}
