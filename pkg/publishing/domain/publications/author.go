package publications

import "omics/pkg/shared/models"

type Author struct {
	ID       models.ID
	Username string
	Name     string
	Lastname string
}

func (a Author) toUser() User {
	return User{
		ID:       a.ID,
		Username: a.Username,
		Name:     a.Name,
		Lastname: a.Lastname,
		Role:     AUTHOR,
	}
}
