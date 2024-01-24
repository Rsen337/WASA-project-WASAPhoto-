package database

import (
	"database/sql"

	"github.com/gofrs/uuid"
)

// returns userId if token is valid and exists or it doesn't
func (db *appdbimpl) DoLogin(token string) (string, error) {

	// Check if the tocken exists and if it does, then return userId
	var userId string
	err := db.c.QueryRow("SELECT id FROM tokens WHERE token = ?", token).Scan(&userId)
	if err == sql.ErrNoRows {
		// Token not found, so create a new user and return a new userId
		u, err := uuid.NewV4()
		if err != nil {
			return "", err
		}
		userId = u.String()

		_, err = db.c.Exec("INSERT INTO tokens (userid, token) VALUES (?, ?)", userId, token)
		if err != nil {
			return "", err
		}

		// *need to create a photos folder and all the associated database for the new user*

		return userId, nil
	} else if err != nil {
		return "", err
	}

	return userId, nil

}
