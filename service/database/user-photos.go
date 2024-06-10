package database

import (
	"errors"
	"time"
)

// GetUserPhotos retrieves the photo IDs for the specified user with pagination.
// It takes the user ID and the page number as parameters and returns a slice of photo IDs and an error.
func (db *appdbimpl) GetUserPhotos(userId string, page int) ([]string, error) {
	// PageSize is the number of photos per page
	const PageSize int = 10

	// Define the SQL query to fetch photo IDs for the specified user with pagination
	offset := (page - 1) * PageSize
	query := "SELECT photoId FROM photos WHERE userId = ? LIMIT ? OFFSET ?"

	// Execute the SQL query
	rows, err := db.c.Query(query, userId, PageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate through the rows and extract photo IDs
	var photoIds []string
	for rows.Next() {
		var photoId string
		if err := rows.Scan(&photoId); err != nil {
			continue
		}
		photoIds = append(photoIds, photoId)
	}

	// Check for any errors during iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return photoIds, nil
}

// UploadPhoto inserts the photo ID into the photo table.
// It takes the photo ID and the user ID as parameters and returns an error.
func (db *appdbimpl) UploadPhoto(photoId string, userId string) error {
	// Insert the photo ID into the photo table
	_, err := db.c.Exec("INSERT INTO photos (photoId, userId, timestamp) VALUES (?, ?, ?)",
		photoId, userId, time.Now().UTC())

	return err
}

// DeletePhoto removes the photo from the database.
// It takes the photo ID as a parameter and returns an error.
func (db *appdbimpl) DeletePhoto(photoId string) error {
	// Remove photo from the database
	stmt, err := db.c.Exec("DELETE FROM photos WHERE photoId = ?", photoId)
	if err != nil {
		return err
	}

	rowsAffected, err := stmt.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("noRows error: photo not found")
	}

	return nil
}
