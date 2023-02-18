package main

func getUserBySession(sessionID string) (int, error) {
	var userID int
	err := db.QueryRow("SELECT user_id FROM sessions WHERE id = ?", sessionID).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}
