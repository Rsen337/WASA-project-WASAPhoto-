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
