package user

import (
	"database/sql"
	"fmt"
	"forum/database"
	e "forum/error"
	"html/template"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func Register(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		e.HandleMethod(w, r)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	validity, errorNum := ValidateUser(username, email, password, db)
	if !validity {
		var errorMessage string
		if errorNum == 1 {
			errorMessage = "You must fill all fields"
		} else if errorNum == 2 {
			errorMessage = "username already exists"
		} else if errorNum == 3 {
			errorMessage = "email already exists"
		} else if errorNum == 5 {
			errorMessage = "Invalid email format"
		} else if errorNum == 7 {
			errorMessage = "Inavlid password format. Your password must be of at least 10 characters length with combination of uppercase, lowercase, numbers and special characters. Example: Passw0rd!2024"
		} else {
			errorMessage = "Sorry, unkown error occured"
		}
		data := errorData{
			Error: errorMessage,
		}
		tmpl, err := template.ParseFiles("web/templates/signup.html")
		if err != nil {
			fmt.Println("Error while parsing template")
			e.HandleInternalError(w, r)
			return
		}
		if err := tmpl.Execute(w, data); err != nil {
			fmt.Println("Error while executing template")
			e.HandleInternalError(w, r)
		}
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		e.HandleInternalError(w, r)
		return
	}
	// Create user record in the database
	email = strings.ToLower(email)
	if _, err := database.InsertUser(db, username, email, string(hashedPassword)); err != nil {
		fmt.Println("Error while inserting user")
		e.HandleInternalError(w, r)
		return
	}

	// Create new user instance
	createNewUser(db, username, email, string(hashedPassword))

	// Redirect to login page after successful registration
	http.Redirect(w, r, "/logIn", http.StatusSeeOther)
}

func createNewUser(db *sql.DB, username, email, password string) User {
	ID, err := database.GetUserID(db, username)
	if err != nil {
		return User{}
	}
	user := User{
		Uid:      ID,
		Email:    email,
		Username: username,
		Password: password,
	}
	return user
}
