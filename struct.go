package main

import "time"

// User represents a registered user of the forum
type User struct {
	ID        int
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time
}

// Post represents a forum post
type Post struct {
	ID         int
	Title      string
	Body       string
	UserID     int64
	CreatedAt  time.Time
	ModifiedAt time.Time
}

// Comment represents a comment on a forum post
type Comment struct {
	ID        int64
	Body      string
	UserID    int64
	PostID    int64
	CreatedAt time.Time
}

// Vote represents a vote on a forum post or comment
type Vote struct {
	ID        int64
	UserID    int64
	PostID    int64
	CommentID int64
	Value     int
}
