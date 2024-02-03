package database

func (db appdbimpl) LikePhoto(photoId string, userId string) error {
	// Prepare the SQL query
	query := "INSERT INTO likes (photoId, userId) VALUES (?, ?)"
	stmt, err := db.c.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the query to insert the new row
	_, err = stmt.Exec(photoId, userId)
	if err != nil {
		return err
	}

	return nil
}

func (db appdbimpl) UnlikePhoto(photoId string, userId string) error {

	// Prepare the SQL query
	query := "DELETE FROM likes WHERE photoId = ? AND userId = ?"
	stmt, err := db.c.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the query to delete the row
	_, err = stmt.Exec(photoId, userId)
	if err != nil {
		return err
	}

	return nil
}
