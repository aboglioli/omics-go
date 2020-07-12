//go:generate mockgen -source $GOFILE -destination mocks/$GOFILE -package mocks
package token

type TokenEncoder interface {
	Encode(id TokenID) (Token, error)
	Decode(t Token) (TokenID, error)
}
