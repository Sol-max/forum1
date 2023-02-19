package main

import (
	"database/sql"
	"html/template"
	"net/http"
	"time"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

var templates = template.Must(template.ParseGlob("templates/*.html"))

type sessionData struct {
	Username string
}

type registrationData struct {
	Username       string
	Email          string
	Password       string
	PasswordRepeat string
}

func main() {
	var err error
	db, err = sql.Open("sqlite3", "forum.db")
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
	)`)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { // home page
		username, err := checkSession(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if username != "" {
			data := &sessionData{Username: username}
			err := templates.ExecuteTemplate(w, "/index.html", data)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
	})
		
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		_, err := checkSession(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if r.Method == "POST" {
			username := r.FormValue("username")
			password := r.FormValue("password")
			user, err := authenticateUser(username, password)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			uuid, err := createSession(user.ID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			expiration := time.Now().Add(time.Hour)
			cookie := &http.Cookie{Name: "session", Value: uuid.String(), Expires: expiration}
			http.SetCookie(w, cookie)
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			err := templates.ExecuteTemplate(w, "login.html", nil)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	})

	http.HandleFunc("/register.html", func(w http.ResponseWriter, r *http.Request) {
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
	
			if registration.Password != registration.PasswordRepeat {
				http.Error(w, "Passwords do not match", http.StatusBadRequest)
				return
			}
	
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registration.Password), bcrypt.DefaultCost)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
	
			createdAt := time.Now().UTC()
			_, err = db.Exec("INSERT INTO users (username, email, password, created_at) VALUES (?, ?, ?, ?)", registration.Username, registration.Email, string(hashedPassword), createdAt)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
	
			user, err := authenticateUser(registration.Username, registration.Password)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
	
			uuid, err := createSession(user.ID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
	
			expiration := time.Now().Add(time.Hour)
			cookie := &http.Cookie{Name: "session", Value: uuid.String(), Expires: expiration}
			http.SetCookie(w, cookie)
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			err := templates.ExecuteTemplate(w, "register.html", nil)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	})
		
	//http.HandleFunc("/", homeHandler)
	http.ListenAndServe(":8080", nil)
}
