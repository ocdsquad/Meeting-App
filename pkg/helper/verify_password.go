package helper

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// VerifyPassword compares a plaintext password with a hashed password.
func VerifyPassword(hashedPassword, plainPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {

		return fmt.Errorf("invalid password")
	}
	return err
}
