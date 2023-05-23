package models

import (
	"time"
)

type Snippet struct {
	ID          int
	UserID      int
	Username    string
	Title       string
	Content     string
	Created     time.Time
	Backend     int
	Frontend    int
	Fullstack   int
	Reactions   int
	snippetID   int
	userID      int
	SumComments int
}

type Session struct {
	ID     int
	UserID int
	Token  string
	Expiry time.Time
}

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Create         time.Time
}

type Comment struct {
	ID        int
	SnippetID int
	UserID    int
	ParentID  int
	Content   string
	Username  string
	Reactions int
	Liked     bool
	Disliked  bool
}
