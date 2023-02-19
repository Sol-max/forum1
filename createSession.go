package main

import "github.com/gofrs/uuid"

func createSession(userID int) (*uuid.UUID, error) {
	// Generate a new UUID for a session
	uuid, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	// Insert into the database generated UUID and the user ID
	_, err = db.Exec("INSERT INTO sessions (uuid, user_id) VALUES (?, ?)", uuid.String(), userID)
	if err != nil {
		return nil, err
	}

	// Return a pointer to new UUID
	return &uuid, nil
}