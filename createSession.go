package main

import "github.com/gofrs/uuid"

func createSession(userID int) (*uuid.UUID, error) {
	// Generate a new UUID for the session
	uuid, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	// Insert the session UUID and the user ID into the sessions table
	_, err = db.Exec("INSERT INTO sessions (uuid, user_id) VALUES (?, ?)", uuid.String(), userID)
	if err != nil {
		return nil, err
	}

	return &uuid, nil
}
