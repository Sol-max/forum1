package forum

import (
	"net/http"

	"github.com/gofrs/uuid"
)
func home(w http.ResponseWriter, r *http.Request) {
	// Get the session cookie from the request
	cookie, err := r.Cookie("session")
	if err != nil {
		// No session cookie found, redirect to login page
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	return
	}

	// Get the session UUID from the cookie value
	uuid, err := uuid.FromString(cookie.Value)
		if err != nil {
			// Invalid session UUID, redirect to login page
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

		// Check if the UUID is valid and retrieve the corresponding user
	user, err := getUserBySession(uuid)
		if err != nil {
			// Failed to retrieve user, redirect to login page
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

		// Render the home page with the user's information
	renderTemplate(w, "home.html", struct {
		Title string
		User  *User
	}{
		Title: "Home",
		User:  user,
	})
}
