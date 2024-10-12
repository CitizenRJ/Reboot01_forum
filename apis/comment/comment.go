package comment

import (
	"database/sql"
	p "forum/apis/post"
	"forum/database"
	e "forum/error"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "modernc.org/sqlite"
)

func CreateNewComment(db *sql.DB, postID int, content string, r *http.Request) Comment {
	userName, err := p.GetUsernameFromSession(r, db)
	if err != nil {
		log.Println(1, err)
		return Comment{}
	}
	userID, err := database.GetUserID(db, userName)
	if err != nil {
		log.Println(2, err)
		return Comment{}
	}

	id, time, err := database.InsertComment(db, postID, userID, content)
	if err != nil {
		log.Println(4, err)
		return Comment{}
	}

	return Comment{
		ID:        int(id),
		PostID:    postID,
		Username:  userName,
		Content:   content,
		CreatedAt: time,
	}
}

func HandleComments(db *sql.DB, content string, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Println(5)
		e.HandleMethod(w, r)
		return
	}

	username, err := p.GetUsernameFromSession(r, db)
	if err != nil || username == "" {
		log.Println(6, err)
		e.HandleStatusForbidden(w, r)
		return
	}

	err = r.ParseForm()
	if err != nil {
		log.Println(7, err)
		e.HandleBadRequest(w, r)
		return
	}
	postID, _ := strconv.Atoi(r.FormValue("postID"))
	// Call CreateNewComment with the username
	CreateNewComment(db, int(postID), content, r)
}

func GetAllComments(db *sql.DB, postID int) ([]Comment, error) {
	query := `SELECT id, post_id, user_id, content, created_at
              FROM comments
              WHERE post_id = ?
              ORDER BY created_at`
	// Prepare the statement
	rows, err := db.Query(query, postID)
	if err != nil {
		log.Println(10, err)
		return nil, err
	}
	defer rows.Close()

	// Slice to hold the results
	var comments []Comment

	// Iterate over the rows
	for rows.Next() {
		var c Comment
		var userID int
		if err != nil {
			log.Println(20, err)
			return nil, err
		}

		// Scan the row into variables
		err = rows.Scan(&c.ID, &c.PostID, &userID, &c.Content, &c.CreatedAt)
		if err != nil {
			log.Println(11, err)
			return nil, err
		}

		loc, _ := time.LoadLocation("Asia/Bahrain") // Set your desired timezone
		c.CreatedAt = c.CreatedAt.In(loc)

		// Retrieve usernames if needed (assuming you have a function for this)
		c.Username, err = database.GetUsernameUsingID(db, userID)
		if err != nil {
			log.Println(12, err)
			return nil, err
		}

		// Append the comment to the slice
		comments = append(comments, c)
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		log.Println(14, err)
		return nil, err
	}

	return comments, nil
}

func GetPostID(db *sql.DB, r *http.Request) int {
	postIDStr := r.FormValue("postID")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		return -1 // Handle error appropriately
	}
	return postID
}
