package users

import (
	"context"
	"omics/pkg/shared/models"
)

type ValidationRepository interface {
	FindByUserID(ctx context.Context, userID models.ID) (*Validation, error)
	Save(ctx context.Context, validation *Validation) error
	Delete(ctx context.Context, userID models.ID) error
}