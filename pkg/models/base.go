package models

import "time"

type ID int

type Base struct {
	ID        ID
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
