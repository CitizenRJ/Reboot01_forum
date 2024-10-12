package user

import (
	"database/sql"
	database "forum/database"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

func checkIfUsernameExists(db *sql.DB, username string) bool {
	usernames, _ := database.GetAllUserNames(db)
	for _, user := range usernames {
		if username == user {
			return true
		}
	}
	return false
}

func checkIfEmaileExists(db *sql.DB, email string) bool {
	emails, _ := database.GetAllUserEmails(db)
	for _, usermail := range emails {
		if strings.EqualFold(usermail, email) {
			return true
		}
	}
	return false
}

func ValidateUser(username, email, password string, db *sql.DB) (bool, int) {
	if username == "" || email == "" || password == "" {
		return false, 1
	}

	if checkIfUsernameExists(db, username) {
		return false, 2
	}

	if !ValidateEmailFormat(email) {
		return false, 5
	}

	if checkIfEmaileExists(db, email) {
		return false, 3
	}

	if !CheckIfPassValid(password) {
		return false, 7
	}

	return true, 0
}

func ValidateLog(userOremail, password string, db *sql.DB) (bool, int) {
	if userOremail == "" || password == "" {
		return false, 1
	}

	if !checkIfUsernameExists(db, userOremail) && !checkIfEmaileExists(db, userOremail) {
		return false, 2
	}

	pass, err := GetUserPassword(db, userOremail)
	if err != nil {
		return false, 4
	}

	err = bcrypt.CompareHashAndPassword([]byte(pass), []byte(password)) //check if both are similar
	if err != nil {
		return false, 3
	}

	return true, 0
}

func GetUserPassword(db *sql.DB, userOremail string) (string, error) {
	isUsername := checkIfUsernameExists(db, userOremail)
	isEmail := checkIfEmaileExists(db, userOremail)

	if isUsername {
		query := `SELECT password FROM users WHERE username = ?`
		var pass string
		err := db.QueryRow(query, userOremail).Scan(&pass)
		if err != nil {
			if err == sql.ErrNoRows {
				// No rows found for the given username
				return "", nil
			}
			return "", err
		}
		return pass, nil
	} else if isEmail {
		email := strings.ToLower(userOremail)
		query := `SELECT password FROM users WHERE email = ?`
		var pass string
		err := db.QueryRow(query, email).Scan(&pass)
		if err != nil {
			if err == sql.ErrNoRows {
				// No rows found for the given username
				return "", nil
			}
			return "", err
		}
		return pass, nil
	}
	return "", nil
}

func ValidateEmailFormat(email string) bool {
	emailRegex := regexp.MustCompile(`^(?i)([a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,})$`)

	// Check if email matches the regexp
	return emailRegex.MatchString(email)
}

func CheckIfPassValid(pass string) bool {
	if len(pass) < 10 {
		return false
	}

	hasUpper := false
	hasLower := false
	hasNumber := false
	hasSpecial := false

	for _, char := range pass {
		if unicode.IsUpper(char) {
			hasUpper = true
		}
		if unicode.IsLower(char) {
			hasLower = true
		}
		if unicode.IsDigit(char) {
			hasNumber = true
		}
		if unicode.IsPunct(char) || unicode.IsSymbol(char) {
			hasSpecial = true
		}
	}
	return hasUpper && hasLower && hasNumber && hasSpecial
}
