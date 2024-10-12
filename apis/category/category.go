package category

import (
	"database/sql"
	p "forum/apis/post"
	"log"
)

func GetPostsByCategory(db *sql.DB, category string) ([]p.Post, error) {
	var posts []p.Post
	query := `
    SELECT p.id, p.user_id, p.title, p.content, p.created_at
    FROM posts p
    JOIN post_categories pc ON p.id = pc.post_id
    JOIN categories c ON pc.category_id = c.id
    WHERE c.name = ?`

	rows, err := db.Query(query, category)
	if err != nil {
		log.Println("Error querying posts by category:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post p.Post
		if err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt); err != nil {
			log.Println("Error scanning post:", err)
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error with rows:", err)
		return nil, err
	}
	return posts, nil
}
