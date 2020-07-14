package publications

import (
	"omics/pkg/shared/models"
	"time"
)

// Publication is an aggregate root
type Publication struct {
	ID          models.ID
	Name        string
	Synopsis    string
	Author      Author
	Pages       []Page
	Statistics  Statistics
	HasContract bool

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
