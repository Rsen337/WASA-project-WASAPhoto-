package database

import "time"

// Photo represents the structure of a photo.
type Photo struct {
	PhotoID     string    `json:"photoID"`     // PhotoID is the unique identifier of the photo.
	UserID      string    `json:"author"`      // UserID is the identifier of the user who uploaded the photo.
	Username	string    `json:"username"`    // Username is the username of the user who uploaded the photo.
	Timestamp   time.Time `json:"timestamp"`   // Timestamp is the time when the photo was uploaded.
	LikesAmount int       `json:"likesAmount"` // LikesAmount is the number of likes the photo has received.
	CommentsAmount    int `json:"commentsAmount"`    // Comments is a list of comments on the photo.
	IsLiked     bool      `json:"isLiked"`     // IsLiked indicates if the photo is liked by the user.
}

// Comment represents the structure of a comment.
type Comment struct {
	CommentID string    `json:"commentID"`   // CommentID is the unique identifier of the comment.
	UserID    string    `json:"author"`      // UserID is the identifier of the user who posted the comment.
	Username	string    `json:"username"`    // Username is the username of the user who posted the comment.
	Content   string    `json:"commentText"` // Content is the text content of the comment.
	Timestamp time.Time `json:"timestamp"`   // Timestamp is the time when the comment was posted.
}

// GetPhoto retrieves the user ID and timestamp of a photo with the given photo ID.
func (db *appdbimpl) GetPhoto(photoID string) (string, time.Time, error) {
	// Prepare the SQL query
	query := "SELECT userId, timestamp FROM photos WHERE photoId = ?"
	stmt, err := db.c.Prepare(query)
	if err != nil {
		return "", time.Now(), err
	}
	defer stmt.Close()

	// Execute the query and scan the result
	var userID string
	var timestamp time.Time
	err = stmt.QueryRow(photoID).Scan(&userID, &timestamp)
	if err != nil {
		return "", time.Now(), err
	}

	return userID, timestamp, nil
}

// GetPhotoLikes retrieves the number of likes for a photo with the given photo ID.
func (db *appdbimpl) GetPhotoLikes(photoID string) (int, error) {
	// Prepare the SQL statement to count likes for the given photo ID
	query := "SELECT COUNT(*) FROM likes WHERE photoId = ?"

	// Execute the query and retrieve the count
	var likesAmount int
	err := db.c.QueryRow(query, photoID).Scan(&likesAmount)
	if err != nil {
		return 0, err
	}
	return likesAmount, nil
}


// GetPhotoCommentsCount retrieves the number of comments for a photo with the given photo ID.
func (db *appdbimpl) GetPhotoCommentsCount(photoID string) (int, error) {
	// Prepare the SQL statement to count comments for the given photo ID
	query := "SELECT COUNT(*) FROM comments WHERE photoId = ?"

	// Execute the query and retrieve the count
	var commentsAmount int
	err := db.c.QueryRow(query, photoID).Scan(&commentsAmount)
	if err != nil {
		return 0, err
	}
	return commentsAmount, nil
}

// GetPhotoComments retrieves the comments for a photo with the given photo ID and page number.
func (db *appdbimpl) GetPhotoComments(photoID string, page int) ([]Comment, error) {
	// PageSize is the number of comments per page
	const PageSize int = 50

	// Define the SQL query to fetch comments for the specified photo with pagination
	offset := (page - 1) * PageSize

	// Prepare the SQL query
	query := "SELECT commentId, userId, commentText, timestamp FROM comments WHERE photoId = ? LIMIT ? OFFSET ?"
	rows, err := db.c.Query(query, photoID, PageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment

	// Iterate over the rows and scan the result into Comment structs
	for rows.Next() {
		var commentID, userID, content string
		var timestamp time.Time

		// Scan the row into variables
		if err := rows.Scan(&commentID, &userID, &content, &timestamp); err != nil {
			return nil, err
		}

		username, err := db.GetUsername(userID)
		if err != nil {
			return nil, err
		}

		// Create a Comment struct and append it to the comments slice
		comment := Comment{
			CommentID: commentID,
			UserID:    userID,
			Username:  username,
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
