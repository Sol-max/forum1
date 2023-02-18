package main

import (
	"net/http"
	"time"

	//"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the user is logged in
	username, err := checkSession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// If the user is logged in, render the home page
	if username != "" {
		data := &sessionData{Username: username}
		err := templates.ExecuteTemplate(w, "home.html", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		// Otherwise, redirect the user to the login page
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
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
		user, err := authenticateUser(username, password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// Create a new session for the user
		uuid, err := createSession(user.ID)
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

func registerHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the user is already logged in
	_, err := checkSession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// If the request method is POST, attempt to register the user
	if r.Method == "POST" {
		// Parse the registration form data
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
		/*errs := validateRegistrationData(data)
		if len(errs) > 0 {
			// If the form data is invalid, render the registration page with error messages
			err := templates.ExecuteTemplate(w, "register.html", map[string]interface{}{
				"Data":  data,
				"Errors": errs,
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}*/

		// Hash the user's password
		hash, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Insert the user into the database
		now := time.Now().UTC().Format(time.RFC3339)
		result, err := db.Exec("INSERT INTO users (username, email, password, created_at) VALUES (?, ?, ?, ?)", data.Username, data.Email, hash, now)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Get the ID of the new user
		userID, err := result.LastInsertId()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Create a new session for the user
		sessionUUID, err := createSession(int(userID))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Set a cookie with the session UUID and redirect the user to the home page
		expiration := time.Now().Add(time.Hour)
		cookie := &http.Cookie{Name: "session", Value: sessionUUID.String(), Expires: expiration}
		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		// Otherwise, render the registration page
		err := templates.ExecuteTemplate(w, "register.html", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}