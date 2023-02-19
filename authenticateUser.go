package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func authenticateUser(username, password string) (*User, error) {
	user, err := getUserByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))// Verify the user's password
	if err != nil {
		return nil, fmt.Errorf("incorrect password")
	}

	return user, nil
}
