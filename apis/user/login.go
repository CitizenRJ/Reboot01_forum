package user

import (
	"database/sql"
	"fmt"
	"forum/database"
	e "forum/error"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/gofrs/uuid"
)

func LogIn(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		e.HandleMethod(w, r)
		return
	}

	err := r.ParseForm()
	if err != nil {
		e.HandleBadRequest(w, r)
		return
	}

	userOremail := r.FormValue("userOremail")
	password := r.FormValue("pass")

	valid, errorNum := ValidateLog(userOremail, password, db)
	if !valid {
		var errorMessage string
		if errorNum == 1 {
			errorMessage = "You must fill all fields"
		} else if errorNum == 2 {
			errorMessage = "username or email does not exists"
		} else if errorNum == 3 {
			errorMessage = "Password is wrong"
		} else {
			errorMessage = "Sorry, unkown error occured"
		}

		data := errorData{
			Error: errorMessage,
		}
		tmpl, err := template.ParseFiles("web/templates/logIn.html")
		if err != nil {
			fmt.Println("Error in parsing template: ", err)
			e.HandleInternalError(w, r)
		}
		if err := tmpl.Execute(w, data); err != nil {
			fmt.Println("Error in executing template: ", err)
			e.HandleInternalError(w, r)
		}
		return
	}

	var userID int
	err = db.QueryRow("SELECT id FROM users WHERE email = ? OR username = ?", strings.ToLower(userOremail), userOremail).Scan(&userID)
	if err != nil {
		e.HandleInternalError(w, r)
		return
	}

	//Check if session exists
	sessionCase, session := checkIfActiveSessionExists(db, userID)
	if sessionCase {
		//end the session
		if err := database.DeleteSession(db, session); err != nil {
			fmt.Println(err)
			e.HandleInternalError(w, r)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:    "session_id",
			Value:   "",
			Path:    "/",
			Expires: time.Now().Add(-time.Hour),
		})
	}

	//New Session must be created
	sessionID := generateSessionID(w, r)              // Implement this function to generate a unique session ID
	expiresAt := time.Now().UTC().Add(24 * time.Hour) // Set session expiration time
	// Retrieve the user ID after validation

	// Insert the session into the database
	if err := database.InsertSession(db, userID, sessionID, expiresAt); err != nil {
		fmt.Println(err)
		e.HandleInternalError(w, r)
	}

	// Set the session cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "session_id",
		Value:   sessionID,
		Path:    "/",
		Expires: expiresAt,
	})
	//Go to postsPage
	http.Redirect(w, r, "/postsPage", http.StatusSeeOther)
}

// generateSessionID generates a unique session ID (implement this according to your needs)
func generateSessionID(w http.ResponseWriter, r *http.Request) string {
	// You can use a UUID or any other method to ensure uniqueness
	new, err := uuid.NewV4()
	if err != nil {
		e.HandleInternalError(w, r)
		return ""
	}
	return new.String() // Replace with actual session ID generation logic
}

func checkIfActiveSessionExists(db *sql.DB, userID int) (bool, int) {
	session, _ := database.GetActiveSessionbyUserID(db, userID)
	return session != -1, session
}
