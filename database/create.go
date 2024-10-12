package database

import (
	"database/sql"
	"fmt"
	"log"
	_ "modernc.org/sqlite"
)

func ConnectToDatabase() *sql.DB {
	db, err := sql.Open("sqlite", "./forum.db")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to the database")
	return db
}

func CreateTables(db *sql.DB) error {
	tableFunctions := []func(*sql.DB) error{
		createUser,
		createCategories,
		createPosts,
		createPostCategories,
		createComments,
		createLikes,
		createSession,
	}

	for _, fn := range tableFunctions {
		if err := fn(db); err != nil {
			return err
		}
	}
	log.Println("All tables created successfully")
	return nil
}

func createUser(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT NOT NULL UNIQUE,
        email TEXT NOT NULL UNIQUE,
        password TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );`
	_, err := db.Exec(query)
	return err
}

func createCategories(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS categories (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL UNIQUE
    );`
	_, err := db.Exec(query)
	return err
}

func createPosts(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS posts (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER,
        title TEXT NOT NULL,
        content TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (user_id) REFERENCES users(id)
    );`
	_, err := db.Exec(query)
	return err
}

func createPostCategories(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS post_categories (
    post_id INTEGER,
    category_id INTEGER,
    PRIMARY KEY (post_id, category_id),
    FOREIGN KEY (post_id) REFERENCES posts(id),
    FOREIGN KEY (category_id) REFERENCES categories(id)
);`
	_, err := db.Exec(query)
	return err
}

func createComments(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS comments (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        post_id INTEGER,
        user_id INTEGER,
        content TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (post_id) REFERENCES posts(id),
        FOREIGN KEY (user_id) REFERENCES users(id),
    );`
	_, err := db.Exec(query)
	return err
}

func createLikes(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS likes (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
        post_id INTEGER,
        comment_id INTEGER,
        is_like BOOLEAN NOT NULL,
        created DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (user_id) REFERENCES users(id),
        FOREIGN KEY (post_id) REFERENCES posts(id),
        FOREIGN KEY (comment_id) REFERENCES comments(id),
        CHECK ((post_id IS NULL AND comment_id IS NOT NULL) OR (post_id IS NOT NULL AND comment_id IS NULL))
    );`
	_, err := db.Exec(query)
	return err
}

func createSession(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS sessions (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER,
        token TEXT NOT NULL UNIQUE,
        expires_at TIMESTAMP NOT NULL,
        FOREIGN KEY (user_id) REFERENCES users(id)
    );`
	_, err := db.Exec(query)
	return err
}
