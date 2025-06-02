package usecases

import (
	"golang.org/x/crypto/bcrypt"
)

type BcryptEncrypter struct{}

func (b *BcryptEncrypter) Encrypt(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}
