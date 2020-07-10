package security

import "fmt"

type fakePasswordHasher struct{}

func FakePasswordHasher() *fakePasswordHasher {
	return &fakePasswordHasher{}
}

func (ph *fakePasswordHasher) Hash(plainPassword string) (string, error) {
	return fmt.Sprintf("#%s#", plainPassword), nil
}

func (ph *fakePasswordHasher) Compare(hashedPassword, plainPassword string) bool {
	return hashedPassword == fmt.Sprintf("#%s#", plainPassword)
}
