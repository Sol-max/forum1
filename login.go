package main

import (
	"net/http"
	"time"
	"golang.org/x/crypto/bcrypt"
)

// Authenticate the user and create a session
func login(w http.ResponseWriter, r *http.Request) {
	// Check if the user is already logged in
		_, err := checkSession(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		return
		}

	// If the request method is POST, attempt to authenticate the user
	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Query the database for the user
		row := db.QueryRow("SELECT id, username, password FROM users WHERE username = ?", username)

		// Check if the user exists
		var userID int
		var dbPassword string
		err := row.Scan(&userID, &username, &dbPassword)
		if err != nil {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		// Compare the user's password to the hashed password stored in the database
		err = bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(password))
		if err != nil {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		// Create a new session for the user
		uuid, err := createSession(int(userID))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Set a cookie with the session UUID and redirect the user to the home page
		expiration := time.Now().Add(time.Hour)
		cookie := &http.Cookie{Name: "session", Value: uuid.String(), Expires: expiration}
		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		// Otherwise, render the login page
		err := templates.ExecuteTemplate(w, "login.html", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
			}
		}
	}