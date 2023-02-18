package main

import (
	"database/sql"
	//"fmt"
	"html/template"
	"net/http"
	"time"
	_ "github.com/mattn/go-sqlite3"

	//"github.com/gofrs/uuid"
	//"github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

// Define a global variable for the database connection
var db *sql.DB

// Define a global variable for the template cache
var templates = template.Must(template.ParseGlob("templates/*.html"))

// Define a struct for the session data
type sessionData struct {
	Username string
}

// Define a struct for the registration form data
type registrationData struct {
	Username       string
	Email          string
	Password       string
	PasswordRepeat string
}

func main() {
	// Initialize the database connection
	var err error
	db, err = sql.Open("sqlite3", "forum.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Create the users table if it doesn't exist
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL
	)`)
	if err != nil {
		panic(err)
	}

	// Serve the home page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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
	})

	// Serve the login page
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
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
	})
	
	
		// Serve the registration page
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
			// Check if the user is already logged in
		_, err := checkSession(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		return
		}

			// If the request method is POST, attempt to register the user
		if r.Method == "POST" {
			username := r.FormValue("username")
			email := r.FormValue("email")
			password := r.FormValue("password")
			passwordRepeat := r.FormValue("password_repeat")

			// Check if the passwords match
			if password != passwordRepeat {
			http.Error(w, "Passwords do not match", http.StatusBadRequest)
			return
			}

			// Hash the password
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Insert the user into the database
			createdAt := time.Now().UTC()
			_, err = db.Exec("INSERT INTO users (username, email, password, created_at) VALUES (?, ?, ?, ?)", username, email, string(hashedPassword), createdAt)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Authenticate the user and create a session
			user, err := authenticateUser(username, password)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

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
			// Otherwise, render the registration page
			err := templates.ExecuteTemplate(w, "register.html", nil)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	})
}