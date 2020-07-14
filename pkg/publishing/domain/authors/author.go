package authors

import "omics/pkg/shared/models"

type Social struct {
	Twitter   string
	Facebook  string
	Instagram string
	LinkedIn  string
}

type Author struct {
	ID       models.ID
	Username string
	Name     string
	Lastname string
	Social   Social
}
