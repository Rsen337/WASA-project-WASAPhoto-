package database

import "database/sql"

// GetFollowings retrieves the list of followings for a given user.
// If page is 0, it returns the whole list without pagination.
func (db appdbimpl) GetFollowings(userId string, page int) ([]User, error) {

	// PageSize is the number of followings per page
	const PageSize int = 50

	// Define the SQL query to fetch followee IDs for the specified user with pagination
	offset := (page - 1) * PageSize

	// Prepare the SQL query
	var rows *sql.Rows
	var err error

	if page == 0 {
		query := "SELECT followeeId FROM followings WHERE followerId = ?"
		rows, err = db.c.Query(query, userId)
		if err != nil {
			return nil, err
		}
	} else {
		query := "SELECT followeeId FROM followings WHERE followerId = ? LIMIT ? OFFSET ?"
		rows, err = db.c.Query(query, userId, PageSize, offset)
		if err != nil {
			return nil, err
		}
	}
	defer rows.Close()

	var followings []User

	// Iterate over the rows and scan the result into User structs
	for rows.Next() {

		var followee User

		// Scan the row into variables
		if err := rows.Scan(&followee.UserID); err != nil {
			return nil, err
		}

		followee.Username, err = db.GetUsername(followee.UserID)
		if err != nil {
			return nil, err
		}

		followings = append(followings, followee)
	}

	// Check for any errors during iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return followings, nil
}

// FollowUser adds a new following relationship between a follower and a followee.
func (db appdbimpl) FollowUser(followerId string, followeeId string) error {
	// Prepare the SQL query
	query := "INSERT INTO followings (followerId, followeeId) VALUES (?, ?)"
	stmt, err := db.c.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the query to insert the new row
	_, err = stmt.Exec(followerId, followeeId)
	if err != nil {
		return err
	}

	return nil
}

// UnfollowUser removes a following relationship between a follower and a followee.
func (db appdbimpl) UnfollowUser(followerId string, followeeId string) error {
	// Prepare the SQL query
	query := "DELETE FROM followings WHERE followerId = ? AND followeeId = ?"
	stmt, err := db.c.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the query to delete the row
	_, err = stmt.Exec(followerId, followeeId)
	if err != nil {
		return err
	}

	return nil
}

// GetFollowers retrieves the list of followers for a given user.
func (db appdbimpl) GetFollowers(userId string, page int) ([]User, error) {

	// PageSize is the number of followers per page
	const PageSize int = 50

	// Define the SQL query to fetch follower IDs for the specified user with pagination
	offset := (page - 1) * PageSize

	// Prepare the SQL query
	query := "SELECT followerId FROM followings WHERE followeeId = ? LIMIT ? OFFSET ?"
	rows, err := db.c.Query(query, userId, PageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var followers []User

	// Iterate over the rows and scan the result into User structs
	for rows.Next() {

		var follower User

		// Scan the row into variables
		if err := rows.Scan(&follower.UserID); err != nil {
			return nil, err
		}

		follower.Username, err = db.GetUsername(follower.UserID)
		if err != nil {
			return nil, err
		}

		followers = append(followers, follower)
	}

	// Check for any errors during iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return followers, nil
}
