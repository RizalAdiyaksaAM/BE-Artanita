package password

import (
	err_util "tugas-akhir/utils/error"
	"golang.org/x/crypto/bcrypt"
)

type PasswordUtil interface {
	HashPassword(password string) (string, error)
	VerifyPassword(password string, hash string) error
}

type passwordUtil struct {}

func NewPasswordUtil() *passwordUtil {
	return &passwordUtil{}
}

func (p *passwordUtil) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err_util.ErrFailedHashingPassword
	}

	return string(hash), nil
}

func (p *passwordUtil) VerifyPassword(password string, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return err_util.ErrPasswordMismatch
	}

	return nil
}