package database

type User struct {
	UserID   string `json:"userId"`
	Username string `json:"username"`
}

func (db appdbimpl) SearchUser(toSearch string, userId string, page int) ([]User, error) {

	// PageSize is the number of photos per page
	const PageSize int = 10

	// Define the SQL query to fetch photo IDs for the specified user with pagination
	offset := (page - 1) * PageSize
	query := "SELECT * FROM users WHERE (username LIKE ?) AND userId NOT IN (SELECT bannerId FROM banned WHERE banneeId = ?) LIMIT ? OFFSET ?"

	// Execute the SQL query
	rows, err := db.c.Query(query, toSearch+"%", userId, PageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate through the rows and extract users
	var users []User
	for rows.Next() {
		var user User
		err = rows.Scan(&user.UserID, &user.Username)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return users, nil
}

type UserProfile struct {
	UserID     string `json:"userId"`
	Username   string `json:"username"`
	Photos     int    `json:"photos"`
	Followers  int    `json:"followers"`
	Followings int    `json:"followings"`
}

func (db appdbimpl) GetUserProfile(userId string) (UserProfile, error) {

	username, err := db.GetUsername(userId)
	if err != nil {
		return UserProfile{}, err
	}

	var photos int
	err = db.c.QueryRow("SELECT COUNT(*) FROM photos WHERE userId = ?", userId).Scan(&photos)
	if err != nil {
		return UserProfile{}, err
	}

	var followers int
	err = db.c.QueryRow("SELECT COUNT(*) FROM followings WHERE followeeId = ?", userId).Scan(&followers)
	if err != nil {
		return UserProfile{}, err
	}

	var followees int
	err = db.c.QueryRow("SELECT COUNT(*) FROM followings WHERE followerId = ?", userId).Scan(&followees)
	if err != nil {
		return UserProfile{}, err
	}

	return UserProfile{
		UserID:     userId,
		Username:   username,
		Photos:     photos,
		Followers:  followers,
		Followings: followees,
	}, nil

}
