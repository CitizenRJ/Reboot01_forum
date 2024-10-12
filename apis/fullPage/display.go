package fullpage

import (
	"database/sql"
	com "forum/apis/category"
	c "forum/apis/comment"
	p "forum/apis/post"
	e "forum/error"
	"html/template"
	"log"
	"net/http"
)

func DisplayHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Check session for username
	username, err := p.GetUsernameFromSession(r, db)
	valid := (err == nil && username != "")

	categoryURL := r.URL.Query().Get("category")

	var postsWithComments []PostWithComments

	if categoryURL == "" {
		// Fetch all posts if no category is specified
		posts, err := p.GetAllPosts(db)
		if err != nil {
			e.HandleInternalError(w, r)
			return
		}

		// Fetch comments for each post
		for _, post := range posts {
			comments, err := c.GetAllComments(db, int(post.ID))
			if err != nil {
				log.Println("Error fetching comments for post ID:", post.ID, "Error:", err)
				e.HandleInternalError(w, r)
				return
			}
			postsWithComments = append(postsWithComments, PostWithComments{Post: post, Comments: comments})
		}
	} else {
		if categoryURL == "Liked" {
			if !valid {
				// Redirect to login page if user is not logged in
				http.Redirect(w, r, "/logIn", http.StatusSeeOther)
				return
			}
			// Fetch liked posts for the logged-in user
			posts, err := GetLikedPosts(db, username)
			if err != nil {
				e.HandleInternalError(w, r)
				return
			}

			// Fetch comments for each liked post
			for _, post := range posts {
				comments, err := c.GetAllComments(db, int(post.ID))
				if err != nil {
					log.Println("Error fetching comments for post ID:", post.ID, "Error:", err)
					e.HandleInternalError(w, r)
					return
				}
				postsWithComments = append(postsWithComments, PostWithComments{Post: post, Comments: comments})
			}
		} else if categoryURL == "UserPosts" {
			posts, err := p.GetUserPost(db, username)
			if err != nil {
				e.HandleInternalError(w, r)
				return
			}
			// Fetch comments for each liked post
			for _, post := range posts {
				comments, err := c.GetAllComments(db, int(post.ID))
				if err != nil {
					log.Println("Error fetching comments for post ID:", post.ID, "Error:", err)
					e.HandleInternalError(w, r)
					return
				}
				postsWithComments = append(postsWithComments, PostWithComments{Post: post, Comments: comments})
			}
		} else {
			// Fetch posts by category
			posts, err := com.GetPostsByCategory(db, categoryURL)
			if err != nil {
				e.HandleInternalError(w, r)
				return
			}

			// Fetch comments for each post in the selected category
			for _, post := range posts {
				comments, err := c.GetAllComments(db, int(post.ID))
				if err != nil {
					log.Println("Error fetching comments for post ID:", post.ID, "Error:", err)
					e.HandleInternalError(w, r)
					return
				}
				postsWithComments = append(postsWithComments, PostWithComments{Post: post, Comments: comments})
			}
		}

	}

	// Parse the template
	tmpl, err := template.ParseFiles("web/templates/postsPage.html")
	if err != nil {
		e.HandleInternalError(w, r)
		return
	}

	// Prepare data for rendering
	data := Page{
		PostsWithComments: postsWithComments,
		Log:               valid,
	}

	// Execute the template with data
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Template execution error: %v", err)
		e.HandleInternalError(w, r)
	}
}

func GetLikedPosts(db *sql.DB, username string) ([]p.Post, error) {
	query := `
        SELECT p.id, p.user_id, p.title, p.content, p.created_at, u.username
        FROM posts p
        JOIN likes l ON p.id = l.post_id
        JOIN users u ON p.user_id = u.id
        JOIN users lu ON l.user_id = lu.id
        WHERE lu.username = ? AND l.is_like = TRUE
        ORDER BY p.created_at DESC
    `
	rows, err := db.Query(query, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []p.Post
	for rows.Next() {
		var post p.Post
		err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt, &post.UserName)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}
