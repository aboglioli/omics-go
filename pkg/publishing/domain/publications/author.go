package publications

import "omics/pkg/shared/models"

type Author struct {
	ID       models.ID
	Username string
	Name     string
	Lastname string
}
