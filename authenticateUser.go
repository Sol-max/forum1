package forum

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func authenticateUser(username, password string) (*User, error) {
	// Look up the user by their username
	user, err := getUserByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	// Verify the user's password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("incorrect password")
	}

	return user, nil
}
