package forum

import (
	"net/http"

	"github.com/gofrs/uuid"
)

// Check if the request has a valid session cookie and return the username associated with the session
func checkSession(r *http.Request) (string, error) {
	cookie, err := r.Cookie("session")
	if err != nil {
		if err == http.ErrNoCookie {
			return "", nil
		}
		return "", err
	}

	// Verify that the session ID is valid and get the corresponding username from the database
	sessionID, err := uuid.FromString(cookie.Value)
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