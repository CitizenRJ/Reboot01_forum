package forum

import (
	"html/template"
	"net/http"
)

func HandleNotFound(w http.ResponseWriter, r *http.Request) {
	// this is where it is recognizing the index.html
	w.WriteHeader(http.StatusNotFound)
	tmpl, err := template.ParseFiles("../forum/web/templates/errorPages/404.html")
	if err != nil {
		HandleInternalError(w, r)
		return
	}
	err = tmpl.Execute(w, "")
	if err != nil {
		HandleInternalError(w, r)
		return
	}

	// http.Error(w, "Page not found", http.StatusNotFound)
}

func HandleBadRequest(w http.ResponseWriter, r *http.Request) {
	// this is where it is recognizing the index.html
	w.WriteHeader(http.StatusBadRequest)
	tmpl, err := template.ParseFiles("../forum/web/templates/errorPages/400.html")
	if err != nil {
		HandleInternalError(w, r)
		return
	}
	err = tmpl.Execute(w, "")
	if err != nil {
		HandleInternalError(w, r)
		return
	}

	// http.Error(w, "Page not found", http.StatusNotFound)
}
func HandleInternalError(w http.ResponseWriter, r *http.Request) {
	// this is where it is recognizing the index.html
	w.WriteHeader(http.StatusInternalServerError)
	tmpl, err := template.ParseFiles("../forum/web/templates/errorPages/500.html")
	if err != nil {
		http.Error(w, "Failed to parse template", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, "")
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}

	// http.Error(w, "Page not found", http.StatusNotFound)
}
func HandleMethod(w http.ResponseWriter, r *http.Request) {
	// this is where it is recognizing the index.html
	w.WriteHeader(http.StatusMethodNotAllowed)
	tmpl, err := template.ParseFiles("../forum/web/templates/errorPages/405.html")
	if err != nil {
		HandleInternalError(w, r)
		return
	}
	err = tmpl.Execute(w, "")
	if err != nil {
		HandleInternalError(w, r)
		return
	}

	// http.Error(w, "Page not found", http.StatusNotFound)
}
func HandleStatusForbidden(w http.ResponseWriter, r *http.Request) {
	// this is where it is recognizing the index.html
	w.WriteHeader(http.StatusForbidden)
	tmpl, err := template.ParseFiles("../forum/web/templates/errorPages/403.html")
	if err != nil {
		HandleInternalError(w, r)
		return
	}
	err = tmpl.Execute(w, "")
	if err != nil {
		HandleInternalError(w, r)
		return
	}

	// http.Error(w, "Page not found", http.StatusNotFound)
}
