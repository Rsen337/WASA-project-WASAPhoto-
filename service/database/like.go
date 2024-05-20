package database

// LikePhoto inserts a new like into the database for the given photo and user.
func (db appdbimpl) LikePhoto(photoID string, userID string) error {
	// Prepare the SQL query
	query := "INSERT INTO likes (photoId, userId) VALUES (?, ?)"
	stmt, err := db.c.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the query to insert the new row
	_, err = stmt.Exec(photoID, userID)
	if err != nil {
		return err
	}

	return nil
}

// UnlikePhoto removes a like from the database for the given photo and user.
func (db appdbimpl) UnlikePhoto(photoID string, userID string) error {

	// Prepare the SQL query
	query := "DELETE FROM likes WHERE photoId = ? AND userId = ?"
	stmt, err := db.c.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the query to delete the row
	_, err = stmt.Exec(photoID, userID)
	if err != nil {
		return err
	}

	return nil
}

// IsLiked checks if a photo is liked by a user.
func (db appdbimpl) IsLiked(photoID string, userID string) (bool, error) {
	// Prepare the SQL query
	query := "SELECT COUNT(*) FROM likes WHERE photoId = ? AND userId = ?"
	stmt, err := db.c.Prepare(query)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	// Execute the query to get the count of rows
	var count int
	err = stmt.QueryRow(photoID, userID).Scan(&count)
	if err != nil {
		return false, err
	}

	// Check if the count is greater than 0
	isLiked := count > 0

	return isLiked, nil
}