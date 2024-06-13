package database

import (
	"github.com/google/uuid"
	"time"
)

// CommentPhoto inserts a comment for a photo into the database.
func (db appdbimpl) CommentPhoto(photoID, userID, commentText string) error {
	// Prepare the SQL statement
	query := "INSERT INTO comments (commentId, photoId, userId, commentText, timestamp) VALUES (?, ?, ?, ?, ?)"
	stmt, err := db.c.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Generate a unique comment ID
	commentID := uuid.New().String()

	// Get the current timestamp
	timestamp := time.Now()

	// Execute the SQL statement to insert the comment
	_, err = stmt.Exec(commentID, photoID, userID, commentText, timestamp)
	if err != nil {
		return err
	}

	return nil
}

// UncommentPhoto deletes a comment from the database.
func (db appdbimpl) UncommentPhoto(commentID string) error {
	// Prepare the SQL statement
	query := "DELETE FROM comments WHERE commentId = ?"
	stmt, err := db.c.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the SQL statement to delete the comment
	_, err = stmt.Exec(commentID)
	if err != nil {
		return err
	}

	return nil
}

// IsPhotoOwner checks if a user is the owner of a photo.
// It returns true if the user is the owner, false otherwise.
// It also returns an error if there was a problem executing the query.
func (db appdbimpl) IsPhotoOwner(photoID, userID string) (bool, error) {

	// Prepare the SQL query
	query := "SELECT EXISTS(SELECT 1 FROM photos WHERE photoId = ? AND userId = ?)"
	stmt, err := db.c.Prepare(query)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	// Execute the SQL query
	var exists bool
	err = stmt.QueryRow(photoID, userID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// IsCommentOwner checks if a user is the owner of a comment.
// It returns true if the user is the owner, false otherwise.
// It also returns an error if there was a problem executing the query.
func (db appdbimpl) IsCommentOwner(commentID, userID string) (bool, error) {
	// Prepare the SQL query
	query := "SELECT EXISTS(SELECT 1 FROM comments WHERE commentId = ? AND userId = ?)"
	stmt, err := db.c.Prepare(query)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	// Execute the SQL query
	var exists bool
	err = stmt.QueryRow(commentID, userID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// CommentExists checks if a comment exists in the database.
// It returns true if the comment exists, false otherwise.
// It also returns an error if there was a problem executing the query.
func (db appdbimpl) CommentExists(commentID string) (bool, error) {
	// Prepare the SQL query to check if the comment exists
	query := "SELECT EXISTS(SELECT 1 FROM comments WHERE commentId = ?)"
	stmt, err := db.c.Prepare(query)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	// Execute the SQL query and scan the result into a boolean variable
	var exists bool
	err = stmt.QueryRow(commentID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
