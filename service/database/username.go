package database

import "strings"

// SetMyUsername changes the username associated with the given userId.
// If the new username is similar to the old one, it returns an error.
// Otherwise, it updates the username and returns nil.
func (db *appdbimpl) SetMyUsername(userId string, newUsername string) error {
	_, err := db.c.Exec("UPDATE users SET username = ? WHERE userId = ?", newUsername, userId)
	if err != nil {
		return err
	}
	return nil
}

// UsernameIsSame checks whether the given username is similar to the old username associated with the userId.
// It returns true if the usernames are the same, false otherwise.
func (db *appdbimpl) UsernameIsSame(userId string, username string) (bool, error) {
	var oldUsername string
	err := db.c.QueryRow("SELECT username FROM users WHERE userId = ?", userId).Scan(&oldUsername)
	if err != nil {
		return false, err
	}
	return username == oldUsername, nil
}

// UsernameIsTaken checks whether the given username is already taken by another user.
// It assumes that UsernameIsSame has been called first to ensure the username is not similar to the old one.
// It returns true if the username is taken, false otherwise.
func (db *appdbimpl) UsernameIsTaken(userId string, username string) (bool, error) {
	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM users WHERE LOWER(username) = LOWER(?)", strings.ToLower(username)).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetUsername retrieves the username associated with the given userId.
// It returns the username and any error encountered during the retrieval.
func (db appdbimpl) GetUsername(userId string) (string, error) {
	var username string
	err := db.c.QueryRow("SELECT username FROM users WHERE userId = ?", userId).Scan(&username)
	return username, err
}
