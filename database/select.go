package database

import (
	"database/sql"
	_ "modernc.org/sqlite"
	"time"
)

func GetAllUserNames(db *sql.DB) ([]string, error) {
	// Define the SQL query to select all rows from the users table
	query := `SELECT username FROM users`

	// Execute the query
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // Ensure rows are closed after processing

	// Slice to hold the results
	var usernames []string
	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			return nil, err
		}
		usernames = append(usernames, username)
	}
	return usernames, nil
}

func GetAllUserEmails(db *sql.DB) ([]string, error) {
	// Define the SQL query to select all rows from the users table
	query := `SELECT email FROM users`

	// Execute the query
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // Ensure rows are closed after processing

	// Slice to hold the results
	var emails []string
	for rows.Next() {
		var email string
		if err := rows.Scan(&email); err != nil {
			return nil, err
		}
		emails = append(emails, email)
	}
	return emails, nil
}

func GetUserID(db *sql.DB, username string) (int, error) {
	query := `SELECT id FROM users WHERE username = ?`
	var id int
	err := db.QueryRow(query, username).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			// No rows found for the given username
			return -1, nil
		}
		return -1, err
	}
	return id, nil
}

func GetUsernameUsingID(db *sql.DB, id int) (string, error) {
	query := `SELECT username FROM users WHERE id = ?`
	var username string
	err := db.QueryRow(query, id).Scan(&username)
	if err != nil {
		if err == sql.ErrNoRows {
			// No rows found for the given username
			return "", nil
		}
		return "", err
	}
	return username, nil
}

func GetPostIDbyUserID(db *sql.DB, userID int) (int, error) {
	query := `SELECT id FROM posts WHERE user_id = ?`
	var id int
	err := db.QueryRow(query, userID).Scan(id)
	if err != nil {
		if err == sql.ErrNoRows {
			// No rows found for the given username
			return -1, nil
		}
		return -1, err
	}
	return id, nil
}

func GetActiveSessionbyUserID(db *sql.DB, userID int) (int, error) {
	query := `SELECT id FROM sessions WHERE user_id = ? AND expires_at > ?`
	var id int
	currentTime := time.Now().UTC()
	err := db.QueryRow(query, userID, currentTime).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			// No rows found or session is expired
			return -1, nil
		}
		return -1, err
	}
	return id, nil
}
