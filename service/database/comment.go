package database

import (
	"github.com/google/uuid"
	"time"
)

func (db appdbimpl) CommentPhoto(photoId string, userId string, commentText string) error {
	// Prepare the SQL statement
	query := "INSERT INTO comments (commentId, photoId, userId, commentText, timestamp) VALUES (?, ?, ?, ?, ?)"
	stmt, err := db.c.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Generate a unique comment ID
	commentId := uuid.New().String()

	// Get the current timestamp
	timestamp := time.Now()

	// Execute the SQL statement to insert the comment
	_, err = stmt.Exec(commentId, photoId, userId, commentText, timestamp)
	if err != nil {
		return err
	}

	return nil
}

func (db appdbimpl) UncommentPhoto(commentId string) error {
	// Prepare the SQL statement
	query := "DELETE FROM comments WHERE commentId = ?"
	stmt, err := db.c.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the SQL statement to delete the comment
	_, err = stmt.Exec(commentId)
	if err != nil {
		return err
	}

	return nil
}

// returns false also if photo doesn't exist
func (db appdbimpl) IsPhotoOwner(photoId string, userId string) (bool, error) {

	// Prepare the SQL query
	query := "SELECT EXISTS(SELECT 1 FROM photos WHERE photoId = ? AND userId = ?)"
	stmt, err := db.c.Prepare(query)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	// Execute the SQL query
	var exists bool
	err = stmt.QueryRow(photoId, userId).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (db appdbimpl) CommentExists(commentId string) (bool, error) {
	// Prepare the SQL query to check if the comment exists
	query := "SELECT EXISTS(SELECT 1 FROM comments WHERE commentId = ?)"
	stmt, err := db.c.Prepare(query)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	// Execute the SQL query and scan the result into a boolean variable
	var exists bool
	err = stmt.QueryRow(commentId).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
