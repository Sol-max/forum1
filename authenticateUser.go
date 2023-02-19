package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// This function authenticates a user if is password matches the password in Database for *User
func authenticateUser(username, password string) (*User, error) {
	// Get the user with the provided username from the database
	user, err := getUserByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	// Compare the provided password with the hashed password in the database
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("incorrect password")
	}

	return user, nil
}
