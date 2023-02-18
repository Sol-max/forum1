package main

import (
	"net/http"
	"text/template"
	"time"

	"github.com/gofrs/uuid"
)

// handles requests to the main page
func homeHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session") // check for the cookie existance
	if err == http.ErrNoCookie {
		uuid, err := uuid.NewV4() // if no cookie, create a new session
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		cookie := &http.Cookie{ // struct creates new session for 1 hour with unique UUID
			Name:    "session",
			Value:   uuid.String(),
			Expires: time.Now().Add(time.Hour),
		}
		http.SetCookie(w, cookie) // send cookie to server
	} else if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	} else {
		// Session cookie present, check if it's valid
		_, err := uuid.FromString(cookie.Value)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)// Invalid session cookie, redirect to login
			return
		}
	}

	// main page
	tmpl, err := template.ParseFiles("template/index.html")
	if err != nil {
		return
	}
	tmpl.Execute(w, db)
}

