package models

import "time"

type ID int

func (id ID) Str() string {
	return string(id)
}

type Base struct {
	ID        ID
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
