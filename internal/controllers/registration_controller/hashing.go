package registration_controller

import "golang.org/x/crypto/bcrypt"

func hashingPassword(password []byte, cost int) ([]byte, error) {
	hashPassword, err := bcrypt.GenerateFromPassword(password, cost)
	if err != nil {
		return nil, err
	}

	return hashPassword, nil
}
