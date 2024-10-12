package post

import (
	"database/sql"
	"fmt"
	"forum/database"
	e "forum/error"
	"log"
	"net/http"
	"time"

	_ "modernc.org/sqlite"
)

func CreateNewPost(db *sql.DB, username, title, content string, categories []string) (Post, error) {
	var userID int
	err := db.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&userID)
	if err != nil {
		return Post{}, err // Handle if user is not found
	}

	id, createdAt, err := database.InsertPost(db, userID, title, content)
	if err != nil {
		return Post{}, err // Handle insert error
	}

	var catId int
	// Check if the category exists
	for _, category := range categories {
		catId, err = GetCategoryID(db, category)
		if err != nil {
			return Post{}, err
		}
		if catId == -1 {
			// If category does not exist, insert it
			catId, err = database.InsertCategory(db, category)
			if err != nil {
				return Post{}, err
			}
		}
		// Insert category for the post
		if err := database.InsertPostCategory(db, int(id), int(catId)); err != nil {
			return Post{}, err
		}
	}

	return Post{
		ID:         id,
		UserName:   username,
		UserID:     userID,
		Title:      title,
		Content:    content,
		Categories: categories,
		CreatedAt:  createdAt,
	}, nil
}

func HandlePostsCreations(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		e.HandleMethod(w, r)
		return
	}

	// Check if the user is logged in
	username, err := GetUsernameFromSession(r, db)
	if err != nil || username == "" {
		e.HandleStatusForbidden(w, r)
		return
	}

	if err := r.ParseForm(); err != nil {
		e.HandleBadRequest(w, r)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")
	category := r.Form["category"] // Change to single category

	// Call CreateNewPost with the username and category
	if _, err := CreateNewPost(db, username, title, content, category); err != nil {
		fmt.Fprintf(w, "Error creating post: %v", err)
		e.HandleInternalError(w, r)
		return
	}
	http.Redirect(w, r, "/postsPage", http.StatusSeeOther)
}

// GetUsernameFromSession retrieves the username from the session.
func GetUsernameFromSession(r *http.Request, db *sql.DB) (string, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return "", err
	}

	// Retrieve user ID from the sessions table using the session ID
	var userID int
	err = db.QueryRow("SELECT user_id FROM sessions WHERE token = ?", cookie.Value).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("session not found")
		}
		return "", err
	}

	// Retrieve the username using the user ID
	username, err := database.GetUsernameUsingID(db, userID)
	if err != nil {
		return "", err
	}

	return username, nil
}

// GetAllPosts retrieves all posts from the database.
func GetAllPosts(db *sql.DB) ([]Post, error) {
	rows, err := db.Query("SELECT id, user_id, title, content, created_at FROM posts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt); err != nil {
			return nil, err
		}

		// Convert post.CreatedAt to the desired timezone
		post.CreatedAt = convertToTimezone(post.CreatedAt, "Asia/Bahrain")

		// Get the username based on user ID
		if post.UserName, err = database.GetUsernameUsingID(db, post.UserID); err != nil {
			log.Println("Error getting username from ID:", err)
		}

		// Append the post to the slice
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

// convertToTimezone converts time to the specified timezone.
func convertToTimezone(t time.Time, location string) time.Time {
	loc, err := time.LoadLocation(location)
	if err != nil {
		log.Println("Error loading location:", err)
		return t // Return the original time if there's an error
	}
	return t.In(loc)
}

func GetCategoryID(db *sql.DB, name string) (int, error) {
	query := `SELECT id FROM categories WHERE name = ?`
	var id int
	err := db.QueryRow(query, name).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return -1, nil // Category does not exist
		}
		return -1, err
	}
	return id, nil
}

// retrieves users posts from the database.
func GetUserPost(db *sql.DB, username string) ([]Post, error) {
	var userID int
	var posts []Post

	// Get the user ID from the username
	err := db.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&userID)
	if err != nil {
		return nil, err // Return the error if the user is not found
	}

	// Query to get posts for the specific user
	query := `
        SELECT p.id, p.user_id, p.title, p.content, p.created_at
        FROM posts p
        WHERE p.user_id = ?
        ORDER BY p.created_at DESC
    `
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Scan each row and append to posts slice
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
