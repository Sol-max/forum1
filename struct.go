package main

import "time"

type User struct {
	ID        int
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time
}

type Post struct {
	ID         int
	Title      string
	Body       string
	UserID     int64
	CreatedAt  time.Time
	ModifiedAt time.Time
}

type Comment struct {
	ID        int64
	Body      string
	UserID    int64
	PostID    int64
	CreatedAt time.Time
}

type Vote struct {
	ID        int64
	UserID    int64
	PostID    int64
	CommentID int64
	Value     int
}
