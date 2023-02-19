package main

import "fmt"

func getUserByUsername(username string) (*User, error) {
	row := db.QueryRow("SELECT id, username, email, password, created_at FROM users WHERE username = ?", username)

	user := &User{}
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	return user, nil
}
