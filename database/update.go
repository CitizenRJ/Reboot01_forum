package database

import (
	"database/sql"
	"time"

	_ "modernc.org/sqlite"
)

// I think not all of them are useful if we won't allow user to change his/her username but updateSession might be useful to extend it.
// FIXME: Unused funciton
func updateUser(db *sql.DB, userID int, newUsername, newEmail, newPassword string) error {
	query := `UPDATE users 
              SET username = ?, email = ?, password = ? 
              WHERE id = ?;`
	_, err := db.Exec(query, newUsername, newEmail, newPassword, userID)
	return err
}

// FIXME: Unused funciton
func updateCategory(db *sql.DB, categoryID int, newName string) error {
	query := `UPDATE categories 
              SET name = ? 
              WHERE id = ?;`
	_, err := db.Exec(query, newName, categoryID)
	return err
}

// FIXME: Unused funciton
func updatePost(db *sql.DB, postID int, newTitle, newContent string) error {
	query := `UPDATE posts 
              SET title = ?, content = ? 
              WHERE id = ?;`
	_, err := db.Exec(query, newTitle, newContent, postID)
	return err
}

// FIXME: Unused funciton
func updatePostCategory(db *sql.DB, postID, oldCategoryID, newCategoryID int) error {
	deleteQuery := `DELETE FROM post_categories 
                    WHERE post_id = ? AND category_id = ?;`
	_, err := db.Exec(deleteQuery, postID, oldCategoryID)
	if err != nil {
		return err
	}

	insertQuery := `INSERT INTO post_categories (post_id, category_id) 
                    VALUES (?, ?);`
	_, err = db.Exec(insertQuery, postID, newCategoryID)
	return err
}

// FIXME: Unused funciton
func updateComment(db *sql.DB, commentID int, newContent string) error {
	query := `UPDATE comments 
              SET content = ? 
              WHERE id = ?;`
	_, err := db.Exec(query, newContent, commentID)
	return err
}

func updateSession(db *sql.DB, sessionID int, newToken string, newExpiresAt time.Time) error {
	query := `UPDATE sessions 
              SET token = ?, expires_at = ? 
              WHERE id = ?;`
	_, err := db.Exec(query, newToken, newExpiresAt, sessionID)
	return err
}
