package web

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	"forum/apis/comment"
	fullpage "forum/apis/fullPage"
	"forum/apis/like"
	"forum/apis/like/repo"
	m "forum/apis/mainP"
	"forum/apis/post"
	"forum/apis/user"
	e "forum/error"
)

func ConnectWeb(db *sql.DB) {
	// // Optionally clear all tables if needed
	// if err := clearAllTables(db); err != nil {
	// 	fmt.Println("Error clearing tables:", err)
	// 	return
	// }

	// Serve static files
	http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("web/"))))
	http.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir("templates/"))))
	http.Handle("/style/", http.StripPrefix("/style/", http.FileServer(http.Dir("style/"))))
	
	// Define routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		m.MainPageHandler(w, r, db)
	})
	http.HandleFunc("/aboutUs", aboutUs)
	http.HandleFunc("/postsPage", func(w http.ResponseWriter, r *http.Request) {
		fullpage.DisplayHandler(db, w, r)
	})
	http.HandleFunc("/createPost", createPost)
	http.HandleFunc("/signUp", signup)
	http.HandleFunc("/logIn", logIn)
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		user.Register(db, w, r)
	})
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		user.LogIn(db, w, r)
	})
	http.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
		post.HandlePostsCreations(db, w, r)
	})

	http.HandleFunc("/comment", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			content := r.FormValue("content")
			// Handle the comment creation in the database
			comment.HandleComments(db, content, w, r)

			// Redirect to the postComment page after successfully adding the comment
			http.Redirect(w, r, "postsPage", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "postsPage", http.StatusSeeOther)
		}
	})

	// likes
	likesRepo := likerepo.NewLikesRepository(db)
	likesService := like.NewLikesService(likesRepo)
	likesController := like.NewLikesController(*likesService)
	http.HandleFunc("/likeDislikePost", likesController.LikeDislikePost)
	http.HandleFunc("/likeDislikeComment", likesController.InteractWithComment)
	http.HandleFunc("/getInteractions", likesController.GetInteractions)

	http.HandleFunc("/logout", LogoutHandler)

	fmt.Println("Listening on: http://localhost:8989/")
	if err := http.ListenAndServe("0.0.0.0:8989", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func aboutUs(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "aboutUs.html")
}

// FIXME: this function is unused
func postsPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "postsPage.html")
}

// FIXME: this function is unused
func postComment(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "postComment.html")
}

func signup(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "signup.html")
}

func logIn(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "logIn.html")
}

func addComment(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "addComment.html")
}

func createPost(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "createPost.html")
}

func getTableNames(db *sql.DB) ([]string, error) {
	query := `SELECT name FROM sqlite_master WHERE type='table'`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, err
		}
		tables = append(tables, tableName)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tables, nil
}

func clearAllTables(db *sql.DB) error {
	tables, err := getTableNames(db)
	if err != nil {
		return err
	}

	for _, table := range tables {
		query := `DELETE FROM ` + table
		_, err := db.Exec(query)
		if err != nil {
			return err
		}
	}

	return nil
}

func renderTemplate(w http.ResponseWriter, tmplName string) {
	tmpl, err := template.ParseFiles("web/templates/" + tmplName)
	if err != nil {
		e.HandleInternalError(w, nil)
		return
	}
	if err := tmpl.Execute(w, nil); err != nil {
		e.HandleInternalError(w, nil)
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Clear the session cookie
	http.SetCookie(w, &http.Cookie{
		Name:   "session_id",
		Value:  "",
		Path:   "/",
		MaxAge: -1, // Expire the cookie
	})

	// Redirect to the homepage or login page
	http.Redirect(w, r, "/logIn", http.StatusSeeOther)
}
