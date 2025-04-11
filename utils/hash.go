package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func ComparePassword(enteredPass, pass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(pass), []byte(enteredPass))
	return err == nil
}
