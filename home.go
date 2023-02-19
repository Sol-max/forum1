package main

import (
	"net/http"
	"time"

	//"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

func homeHandler(w http.ResponseWriter, r *http.Request) { // Check if the user is logged in
	username, err := checkSession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if username != "" { // If the user is logged in, return to the home page
		data := &sessionData{Username: username}
		err := templates.ExecuteTemplate(w, "/index.html", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther) // Return the user to the login page
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) { // Check if the user is already logged in
	_, err := checkSession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Method == "POST" { // If the request method is POST, attempt to authenticate the user
		username := r.FormValue("username")
		password := r.FormValue("password")
		user, err := authenticateUser(username, password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		uuid, err := createSession(user.ID) // Create a new session for the user
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		expiration := time.Now().Add(time.Hour) // Set a cookie with the session UUID and return the user to the home page
		cookie := &http.Cookie{Name: "session", Value: uuid.String(), Expires: expiration}
		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		err := templates.ExecuteTemplate(w, "login.html", nil) // Return to the login page
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) { // Check if the user is already logged in
	_, err := checkSession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Method == "POST" { // If the request method is POST, attempt to register the user
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Validate the form data
		data := registrationData{
			Username:       r.PostFormValue("username"),
			Email:          r.PostFormValue("email"),
			Password:       r.PostFormValue("password"),
			PasswordRepeat: r.PostFormValue("password_repeat"),
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost) // Save the password
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		now := time.Now().UTC().Format(time.RFC3339) // Insert the user into the database
		result, err := db.Exec("INSERT INTO users (username, email, password, created_at) VALUES (?, ?, ?, ?)", data.Username, data.Email, hash, now)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		userID, err := result.LastInsertId() // Get the ID of the new user
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		sessionUUID, err := createSession(int(userID)) // Create a new session for the user
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		expiration := time.Now().Add(time.Hour) // Set a cookie with the session UUID and redirect the user to the home page
		cookie := &http.Cookie{Name: "session", Value: sessionUUID.String(), Expires: expiration}
		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {

		err := templates.ExecuteTemplate(w, "register.html", nil) // Return to the registration page
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
