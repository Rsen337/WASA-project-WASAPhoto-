package database

// GetMyStream retrieves the stream of photos for a given user ID and page number.
// It returns a slice of Photo objects and an error, if any.
func (db *appdbimpl) GetMyStream(userID string, page int) ([]Photo, error) {
	// PageSize defines the number of photos per page
	const PageSize = 5

	offset := (page - 1) * PageSize

	query := `
		SELECT * FROM photos
		WHERE userID IN (
			SELECT followeeID FROM followings WHERE followerID = ?
		)
		ORDER BY timestamp DESC
		LIMIT ? OFFSET ?;
	`

	rows, err := db.c.Query(query, userID, PageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Retrieve all the photos from the followings
	var photos []Photo

	for rows.Next() {
		var photo Photo
		err = rows.Scan(&photo.PhotoID, &photo.UserID, &photo.Timestamp)
		if err != nil {
			return nil, err
		}

		// Retrieve comments for the photo
		photo.CommentsAmount, err = db.GetPhotoCommentsCount(photo.PhotoID)
		if err != nil {
			return nil, err
		}

		// Retrieve likes amount for the photo
		photo.LikesAmount, err = db.GetPhotoLikes(photo.PhotoID)
		if err != nil {
			return nil, err
		}

		photos = append(photos, photo)
	}

	// Check for any errors during iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return photos, nil
}
