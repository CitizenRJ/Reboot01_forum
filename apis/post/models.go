package post

import (
	"time"
)

type Post struct {
	ID        int64
	UserName  string
	UserID    int
	Title     string
	Content   string
	Categories  []string // Store multiple categories as a slice of strings
	CreatedAt time.Time
}

type Page struct {
	PagePost []Post
	Log      bool
}
