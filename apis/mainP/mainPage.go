package mainP

import (
	"database/sql"
	post "forum/apis/post"
	e "forum/error"
	"html/template"
	"log"
	"net/http"
)

func MainPageHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	path := r.URL.Path
	if path != "/" && path != "/logIn"&&path != "/signup"&&path != "/aboutUs"&&path != "/createPost"&&path != "/postsPage"{
		e.HandleBadRequest(w,r)
		return

	}
	username, err := post.GetUsernameFromSession(r, db)
	valid := (err == nil && username != "")
	tmpl, err := template.ParseFiles("web/templates/mainPage.html")
	if err != nil {
		e.HandleInternalError(w, r)
		return
	}
	data := Page{
		Log: valid,
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		e.HandleInternalError(w, r)
		return
	}

}
