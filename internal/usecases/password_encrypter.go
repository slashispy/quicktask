package usecases

type PasswordEncrypter interface {
	Encrypt(password string) (string, error)
}
