package repository

import (
	"errors"
)

var (
	ErrOperatorAlreadyExists = errors.New("operator already exists")
	ErrActiveDialogExists    = errors.New("active dialog already exists")
	ErrArchivedDialogExists  = errors.New("archived dialog already exists")
	ErrQueuedDialogExists    = errors.New("queued dialog already exists")
	ErrClientAlreadyExists   = errors.New("client already exists")
)
