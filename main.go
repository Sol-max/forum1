package main

import (
	"database/sql"
	"net/http"
)

// Define a global variable to store the database connection
var db *sql.DB

// Define a struct to represent the main page
type IndexPage struct {
	Title       string
	Description string
	Posts       []Post
}

// Define a struct to represent a post
type Post struct {
	Title       string
	Content     string
	Category    string
	AuthorEmail string
	Likes       int
	Dislikes    int
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.ListenAndServe(":8080", nil)
}