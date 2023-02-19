package main

import (
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func login(w http.ResponseWriter, r *http.Request) { // Authenticate the user and create a session
	_, err := checkSession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Method == "POST" { // If the request method is POST, attempt to authenticate the user
		username := r.FormValue("username")
		password := r.FormValue("password")

		row := db.QueryRow("SELECT id, username, password FROM users WHERE username = ?", username) // fill the user data
		var userID int
		var dbPassword string
		err := row.Scan(&userID, &username, &dbPassword) // Check if the user exists
		if err != nil {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}
		err = bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(password)) // Compare the password to the saved passwords stored in the database
		if err != nil {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		uuid, err := createSession(int(userID)) // Create a new session for the user
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		expiration := time.Now().Add(time.Hour) // Set a cookie with the session UUID and return the user to the home page
		cookie := &http.Cookie{Name: "session", Value: uuid.String(), Expires: expiration}
		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/index", http.StatusSeeOther)
	} else {
		err := templates.ExecuteTemplate(w, "login.html", nil) // Return to the login page
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
