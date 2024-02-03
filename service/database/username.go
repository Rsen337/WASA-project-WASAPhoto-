package database

// function that chagnes the username to new one. If the username is similar to the old one, then it returns false, otherwise true
func (db *appdbimpl) SetMyUsername(userId string, newUsername string) error {

	// update username associated to the userId
	_, err := db.c.Exec("UPDATE users SET username = ? WHERE userId = ?", newUsername, userId)
	if err != nil {
		return err
	}

	return nil
}

// checks whether the username is similar to the old one
func (db *appdbimpl) UsernameIsSame(userId string, username string) (bool, error) {
	var oldUsername string
	err := db.c.QueryRow("SELECT username FROM users WHERE userId = ?", userId).Scan(&oldUsername)
	if err != nil {
		return false, err
	}

	return username == oldUsername, nil
}

// checks whether the username is already taken by some other user
// it assumes that we ran UsernameIsSame first
func (db *appdbimpl) UsernameIsTaken(userId string, username string) (bool, error) {
	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", username).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (db appdbimpl) GetUsername(userId string) (string, error) {
	var username string
	err := db.c.QueryRow("SELECT username FROM users WHERE userId = ?", userId).Scan(&username)
	return username, err
}
