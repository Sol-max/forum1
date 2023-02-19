package main

import "github.com/gofrs/uuid"

func createSession(userID int) (*uuid.UUID, error) {
	uuid, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("INSERT INTO sessions (uuid, user_id) VALUES (?, ?)", uuid.String(), userID)
	if err != nil {
		return nil, err
	}

	return &uuid, nil
}
