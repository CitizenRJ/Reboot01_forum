package comment

import "time"

type Comment struct {
	ID        int
	PostID    int
	Username  string
	Content   string
	CreatedAt time.Time
}

type Page struct {
	CommentPage []Comment
	Log         bool
}
