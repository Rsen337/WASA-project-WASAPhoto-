package database

// returns false also if banner == bannee
func (db appdbimpl) IsBanned(bannerId string, banneeId string) (bool, error) {

	if bannerId == banneeId {
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
	err = stmt.QueryRow(bannerId, banneeId).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (db appdbimpl) GetBannedUsers(userId string, page int) ([]User, error) {

	// PageSize is the number of photos per page
	const PageSize int = 50

	// Define the SQL query to fetch photo IDs for the specified user with pagination
	offset := (page - 1) * PageSize

	// Prepare the SQL query
	query := "SELECT banneeId FROM banned WHERE bannerId = ? LIMIT ? OFFSET ?"
	rows, err := db.c.Query(query, userId, PageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bannedUsers []User

	// Iterate over the rows and scan the result into Comment structs
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

func (db appdbimpl) BanUser(bannerId string, banneeId string) error {
	// Prepare the SQL query
	query := "INSERT INTO banned (bannerId, banneeId) VALUES (?, ?)"
	stmt, err := db.c.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the query to insert the new row
	_, err = stmt.Exec(bannerId, banneeId)
	if err != nil {
		return err
	}

	return nil
}

func (db appdbimpl) UnbanUser(bannerId string, banneeId string) error {
	// Prepare the SQL query
	query := "DELETE FROM banned WHERE bannerId = ? AND banneeId = ?"
	stmt, err := db.c.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the query to delete the row
	_, err = stmt.Exec(bannerId, banneeId)
	if err != nil {
		return err
	}

	return nil
}
