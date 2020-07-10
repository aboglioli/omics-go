package token

type TokenEncoder interface {
	Encode(tokenID TokenID) (Token, error)
	Decode(token Token) (TokenID, error)
}
