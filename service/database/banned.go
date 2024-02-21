package database

// IsBanned checks if a user is banned by another user.
// It returns false if the banner and bannee IDs are the same.
func (db appdbimpl) IsBanned(bannerID string, banneeID string) (bool, error) {

	if bannerID == banneeID {
		return false, nil
	}

	// Prepare the SQL query
	query := "SELECT EXISTS(SELECT 1 FROM banned WHERE bannerId = ? AND banneeId = ?)"
	stmt, err := db.c.Prepare(query)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	// Execute the query and scan the result
	var exists bool
	err = stmt.QueryRow(bannerID, banneeID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// GetBannedUsers retrieves a list of banned users for a given user ID with pagination.
func (db appdbimpl) GetBannedUsers(userID string, page int) ([]User, error) {

	// PageSize is the number of users per page
	const PageSize int = 50

	// Define the SQL query to fetch banned user IDs for the specified user with pagination
	offset := (page - 1) * PageSize

	// Prepare the SQL query
	query := "SELECT banneeId FROM banned WHERE bannerId = ? LIMIT ? OFFSET ?"
	rows, err := db.c.Query(query, userID, PageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bannedUsers []User

	// Iterate over the rows and scan the result into User structs
	for rows.Next() {

		var bannee User

		// Scan the row into variables
		if err := rows.Scan(&bannee.UserID); err != nil {
			return nil, err
		}

		bannee.Username, err = db.GetUsername(bannee.UserID)
		if err != nil {
			return []User{}, err
		}

		bannedUsers = append(bannedUsers, bannee)
	}

	// Check for any errors during iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return bannedUsers, nil
}

// BanUser bans a user by adding a new row to the "banned" table.
func (db appdbimpl) BanUser(bannerID string, banneeID string) error {
	// Prepare the SQL query
	query := "INSERT INTO banned (bannerId, banneeId) VALUES (?, ?)"
	stmt, err := db.c.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the query to insert the new row
	_, err = stmt.Exec(bannerID, banneeID)
	if err != nil {
		return err
	}

	return nil
}

// UnbanUser removes a user from the "banned" table.
func (db appdbimpl) UnbanUser(bannerID string, banneeID string) error {
	// Prepare the SQL query
	query := "DELETE FROM banned WHERE bannerId = ? AND banneeId = ?"
	stmt, err := db.c.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the query to delete the row
	_, err = stmt.Exec(bannerID, banneeID)
	if err != nil {
		return err
	}

	return nil
}
