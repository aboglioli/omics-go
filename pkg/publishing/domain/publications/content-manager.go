package publications

import "omics/pkg/shared/models"

type ContentManager struct {
	id       models.ID
	username string
	name     string
	lastname string
}

func (cm *ContentManager) ID() models.ID {
	return cm.id
}
