package fullpage

import (
	c "forum/apis/comment"
	p "forum/apis/post"
)

type PostWithComments struct {
	Post     p.Post
	Comments []c.Comment // Assuming a Comment struct exists in your comment package
}

type Page struct {
	PostsWithComments []PostWithComments
	Log               bool
}
