package models

import (
	"strconv"
	"time"
)

type ID int

func (id ID) String() string {
	return strconv.Itoa(int(id))
}

type Base struct {
	ID        ID
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
