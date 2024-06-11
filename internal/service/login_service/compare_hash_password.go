package login_service

import "golang.org/x/crypto/bcrypt"

func ComparePasswords(hashed string, plain []byte) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), plain)
	return err == nil
}
