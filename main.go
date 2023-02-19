package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var (
	db *sql.DB
	// HTML templates in use
	templates = template.Must(template.ParseGlob("templates/*.html"))
)

// user data
type sessionData struct {
	Username string
}

// data struct for registration
type registrationData struct {
	Username       string
	Email          string
	Password       string
	PasswordRepeat string
}

func main() {
	var err error
	db, err = sql.Open("sqlite3", "forum.db") // connecting to database
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users ( 
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL
	)`) // creates user table
	if err != nil {
		panic(err)
	}

	// handlers for each page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { // home page
		username, err := checkSession(r) // Check user session.
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if username != "" { // If session is valid, run home page to active user
			data := &sessionData{Username: username}
			err := templates.ExecuteTemplate(w, "/index.html", data)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else { // If session is not valid, return to the login page.
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) { // Login page
		_, err := checkSession(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if r.Method == "POST" { // If the request method is POST, authenticate the user.
			username := r.FormValue("username")
			password := r.FormValue("password")
			user, err := authenticateUser(username, password)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			// If the user is authenticated, create a new session and set a session cookie.
			uuid, err := createSession(user.ID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			expiration := time.Now().Add(time.Hour)
			cookie := &http.Cookie{Name: "session", Value: uuid.String(), Expires: expiration}
			http.SetCookie(w, cookie)
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else { // If the request method is GET, return to the login page.
			err := templates.ExecuteTemplate(w, "login.html", nil)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	})

	http.HandleFunc("/register.html", func(w http.ResponseWriter, r *http.Request) { // Register page
		_, err := checkSession(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if r.Method == "POST" {
			registration := registrationData{
				Username:       r.FormValue("username"),
				Email:          r.FormValue("email"),
				Password:       r.FormValue("password"),
				PasswordRepeat: r.FormValue("password_repeat"),
			}
			// Check if the passwords match
			if registration.Password != registration.PasswordRepeat {
				http.Error(w, "Passwords do not match", http.StatusBadRequest)
				return
			}
			// Hashes the user password
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registration.Password), bcrypt.DefaultCost)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			createdAt := time.Now().UTC() // Put the new user into the database
			_, err = db.Exec("INSERT INTO users (username, email, password, created_at) VALUES (?, ?, ?, ?)", registration.Username, registration.Email, string(hashedPassword), createdAt)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// Authenticate the newly registered user
			user, err := authenticateUser(registration.Username, registration.Password)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// Create a new session for the user
			uuid, err := createSession(user.ID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// Create a cookie for the sessionI
			expiration := time.Now().Add(time.Hour)
			cookie := &http.Cookie{Name: "session", Value: uuid.String(), Expires: expiration}
			http.SetCookie(w, cookie)
			http.Redirect(w, r, "/", http.StatusSeeOther) // Then return user to the Home page
		} else { // For GET request return to the Register page
			err := templates.ExecuteTemplate(w, "register.html", nil)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	})

	// http.HandleFunc("/", homeHandler)
	log.Println("Server listening on :8080...")
	http.ListenAndServe(":8080", nil)
	http.HandleFunc("/login", login)
	// http.HandleFunc("/",homeHandler)
}
