package token

import (
	"context"
	"omics/pkg/shared/models"

	"github.com/google/uuid"
)

const (
	TOKEN_KEY = "auth_token"
)

// TokenID
type TokenID string

func NewTokenID() TokenID {
	uuid := uuid.New().String()
	return TokenID(uuid)
}

// Token
type Token string

func (t Token) String() string {
	return string(t)
}

func (t Token) ToContext(parent context.Context) context.Context {
	return context.WithValue(parent, TOKEN_KEY, t)
}

// Data
type Data map[string]string

func NewData(userID models.ID) Data {
	return Data{
		"user_id": string(userID),
	}
}

func (d Data) UserID() (models.ID, error) {
	if userID, ok := d["user_id"]; ok {
		return models.ID(userID), nil
	}
	return "", Err.Code("extracting_user_id_from_data")
}

func FromContext(ctx context.Context) (Token, error) {
	if token, ok := ctx.Value(TOKEN_KEY).(Token); ok {
		return token, nil
	}
	return "", Err.Code("token_from_context")
}
