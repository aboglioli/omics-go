package publications

import "omics/pkg/shared/models"

const (
	CONTENT_MANAGER = "content-manager"
	AUTHOR          = "author"
)

type User struct {
	ID       models.ID
	Username string
	Name     string
	Lastname string
	Role     string
}
