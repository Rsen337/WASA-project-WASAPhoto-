package database

import "time"

// Photo represents the structure of a photo.
type Photo struct {
	PhotoID     string    `json:"photoID"`
	UserID      string    `json:"author"`
	Timestamp   time.Time `json:"timestamp"`
	LikesAmount int       `json:"likesAmount"`
	Comments    []Comment `json:"comments"`
}

// Comment represents the structure of a comment.
type Comment struct {
	CommentID string    `json:"commentID"`
	UserID    string    `json:"author"`
	Content   string    `json:"commentText"`
	Timestamp time.Time `json:"timestamp"`
}

// Getauthor is an example that shows you how to query data
func (db *appdbimpl) GetPhoto(photoId string) (string, time.Time, error) {

	// Prepare the SQL query
	query := "SELECT userId, timestamp FROM photos WHERE photoId = ?"
	stmt, err := db.c.Prepare(query)
	if err != nil {
		return "", time.Now(), err
	}
	defer stmt.Close()

	// Execute the query and scan the result
	var userId string
	var timestamp time.Time
	err = stmt.QueryRow(photoId).Scan(&userId, &timestamp)
	if err != nil {
		return "", time.Now(), err
	}

	return userId, timestamp, nil
}

func (db *appdbimpl) GetPhotoLikes(photoId string) (int, error) {
	// Prepare the SQL statement to count likes for the given photoId
	query := "SELECT COUNT(*) FROM likes WHERE photoId = ?"

	// Execute the query and retrieve the count
	var likesAmount int
	err := db.c.QueryRow(query, photoId).Scan(&likesAmount)
	if err != nil {
		return 0, err
	}
	return likesAmount, nil
}

func (db *appdbimpl) GetPhotoComments(photoId string) ([]Comment, error) {
	var comments []Comment

	// Prepare the SQL query
	query := "SELECT commentId, userId, commentText, timestamp FROM comments WHERE photoId = ?"
	rows, err := db.c.Query(query, photoId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows and scan the result into Comment structs
	for rows.Next() {
		var commentID, userID, content string
		var timestamp time.Time

		// Scan the row into variables
		if err := rows.Scan(&commentID, &userID, &content, &timestamp); err != nil {
			return nil, err
		}

		// Create a Comment struct and append it to the comments slice
		comment := Comment{
			CommentID: commentID,
			UserID:    userID,
			Content:   content,
			Timestamp: timestamp,
		}
		comments = append(comments, comment)
	}

	// Check for any errors during iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}
