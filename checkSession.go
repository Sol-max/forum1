package main

import (
	"net/http"

	"github.com/gofrs/uuid"
)

func checkSession(r *http.Request) (string, error) {// Check if the request has a valid session cookie and return the username
	cookie, err := r.Cookie("session")
	if err != nil {
		if err == http.ErrNoCookie {
			return "", nil
		}
		return "", err
	}

	sessionID, err := uuid.FromString(cookie.Value)	// Verify that the session ID is valid
	if err != nil {
		return "", err
	}

	var username string
	err = db.QueryRow("SELECT username FROM sessions WHERE id=?", sessionID.String()).Scan(&username)
	if err != nil {
		return "", err
	}

	return username, nil
}