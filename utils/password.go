package utils

import "golang.org/x/crypto/bcrypt"

func HashingPassword(password string) string {
	hashedByte, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return ""
	}
	return string(hashedByte)
}

func CheckPasswordHash(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false
	}
	return true
}
