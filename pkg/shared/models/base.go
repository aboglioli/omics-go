package models

import (
	"github.com/google/uuid"
)

// Entity must have an ID
type Entity interface {
	ID() ID
	Equals(entity interface{}) bool
}

// ID is an unique identifier used in each bounded context
type ID string

func NewID() ID {
	return ID(uuid.New().String())
}

func (id ID) Equals(otherID ID) bool {
	return id == otherID
}

func (id ID) String() string {
	return string(id)
}
