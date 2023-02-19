package main

import "fmt"

func getUserByUsername(username string) (*User, error) {
	// Get from the database
	row := db.QueryRow("SELECT id, username, email, password, created_at FROM users WHERE username = ?", username)

	// Create a new User struct
	user := &User{}

	// Scan the results row into the User struct
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		// If there was an error scanning the result row, return an error with a helpful message
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	// Return the User struct and no error
	return user, nil
}
