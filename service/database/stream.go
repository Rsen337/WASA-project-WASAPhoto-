package database

import "sort"

// SetName is an example that shows you how to execute insert/update
func (db *appdbimpl) GetMyStream(userId string, page int) ([]Photo, error) {

	// PageSize is the number of photos per page
	const PageSize int = 50

	followings, err := db.GetFollowings(userId, 0)
	if err != nil {
		return nil, err
	}

	query := "SELECT * FROM photos WHERE userId = ?"

	// get all the photos from all the followings
	var photos []Photo
	for _, following := range followings {
		rows, err := db.c.Query(query, following.UserID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var photo Photo
			err = rows.Scan(&photo.PhotoID, &photo.UserID, &photo.Timestamp)
			if err != nil {
				return nil, err
			}
			photo.Comments, err = db.GetPhotoComments(photo.PhotoID)
			if err != nil {
				return []Photo{}, err
			}

			photo.LikesAmount, err = db.GetPhotoLikes(photo.PhotoID)
			if err != nil {
				return []Photo{}, err
			}

			photos = append(photos, photo)
		}
	}

	// sort the photo by chronological order
	sort.Slice(photos, func(i, j int) bool {
		return photos[i].Timestamp.After(photos[j].Timestamp)
	})

	// Calculate start and end indices
	start := (page - 1) * PageSize
	end := start + PageSize

	// Check if end index is within the range of the slice
	if end > len(photos) {
		end = len(photos)
	}

	// Get the photos for the current page
	pagePhotos := photos[start:end]

	return pagePhotos, nil
}
