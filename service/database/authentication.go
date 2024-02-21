package database

import (
	"database/sql"

	"github.com/gofrs/uuid"
)

// DoLogin checks if the token is valid and exists in the database.
// If the token exists, it returns the corresponding userId.
// If the token doesn't exist, it creates a new user and returns a new userId.
func (db *appdbimpl) DoLogin(token string) (string, bool, error) {

	var userId string

	// Check if the token exists in the database
	err := db.c.QueryRow("SELECT userId FROM users WHERE username = ?", token).Scan(&userId)
	if err == sql.ErrNoRows {
		// Token not found, so create a new user and return a new userId
		u, err := uuid.NewV4()
		if err != nil {
			return "", false, err
		}
		userId = u.String()

		_, err = db.c.Exec("INSERT INTO users (userId, username) VALUES (?, ?)", userId, token)
		if err != nil {
			return "", false, err
		}

		return userId, false, nil
	} else if err != nil {
		return "", false, err
	}

	return userId, true, nil
}

// UserExists checks if a user with the given userId exists in the database.
func (db *appdbimpl) UserExists(userId string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS (SELECT 1 FROM users WHERE userId = ?)"
	err := db.c.QueryRow(query, userId).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
