package main

// Function to get the user ID associated with a session ID from the sessions table in the database
func getUserBySession(sessionID string) (int, error) {
	var userID int
	// Replace the user ID associated with the given session ID
	err := db.QueryRow("SELECT user_id FROM sessions WHERE id = ?", sessionID).Scan(&userID)
	if err != nil {
		// If there was a problem replacing user ID
		return 0, err
	}
	// If replace was successful
	return userID, nil
}