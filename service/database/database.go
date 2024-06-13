/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {

	// Register or log in user
	DoLogin(token string) (string, bool, error)
	UserExists(userId string) (bool, error)

	// Username
	SetMyUsername(userId string, newUsername string) error
	UsernameIsSame(userId string, username string) (bool, error)
	UsernameIsTaken(userId string, username string) (bool, error)
	GetUsername(userId string) (string, error)

	// My photos
	GetUserPhotos(userId string, page int) ([]string, error)
	UploadPhoto(photoId string, userId string) error
	DeletePhoto(photId string) error

	// Photos
	GetPhoto(photoId string) (string, time.Time, error)
	GetPhotoLikes(photoId string) (int, error)
	GetPhotoComments(photoId string, page int) ([]Comment, error)

	// Followers
	GetFollowers(userId string, page int) ([]User, error)

	// Followings
	GetFollowings(userId string, page int) ([]User, error)
	FollowUser(followerId string, followeeId string) error
	UnfollowUser(followerId string, followeeId string) error
	IsFollowed(followerId string, followeeId string) (bool, error)

	// Banned
	IsBanned(bannerId string, banneeId string) (bool, error)
	GetBannedUsers(userId string, page int) ([]User, error)
	BanUser(followerId string, followeeId string) error
	UnbanUser(followerId string, followeeId string) error

	// Search users
	SearchUser(toSearh string, userId string, page int) ([]User, error)
	GetUserProfile(userId string, authUser string) (UserProfile, error)

	// Stream
	GetMyStream(userId string, page int) ([]Photo, error)

	// Like
	LikePhoto(photoId string, userId string) error
	UnlikePhoto(photoId string, userId string) error
	IsLiked(photoID string, userID string) (bool, error)

	// Comment
	CommentPhoto(photoId string, userId string, commentText string) error
	UncommentPhoto(commentId string) error
	IsPhotoOwner(photoId string, userId string) (bool, error)
	CommentExists(commentId string) (bool, error)
	GetPhotoCommentsCount(photoID string) (int, error)
	IsCommentOwner(commentId string, userId string) (bool, error)

	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Enable foreign keys
	_, errPramga := db.Exec(`PRAGMA foreign_keys= ON`)
	if errPramga != nil {
		return nil, fmt.Errorf("error setting pragmas: %w", errPramga)
	}

	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	_, err := db.Exec(`DROP TABLE IF EXISTS example_table`)
	if err != nil {
		return nil, err
	}
	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='example_table';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {

		// define all the necessary tables
		sqlStmts := [7]string{
			`CREATE TABLE IF NOT EXISTS example_table (
				id INTEGER NOT NULL PRIMARY KEY, name TEXT
			);`,
			`CREATE TABLE IF NOT EXISTS users (
				userId TEXT PRIMARY KEY,
				username TEXT UNIQUE
			);`,
			`CREATE TABLE IF NOT EXISTS photos (
				photoId TEXT PRIMARY KEY,
				userId TEXT NOT NULL,
				timestamp DATETIME NOT NULL,
				FOREIGN KEY(userId) REFERENCES users (userId) ON DELETE CASCADE
			);`,
			`CREATE TABLE IF NOT EXISTS likes (
				photoId TEXT NOT NULL,
				userId TEXT NOT NULL,
				PRIMARY KEY (photoId,userId),
				FOREIGN KEY(photoId) REFERENCES photos (photoId) ON DELETE CASCADE
			);`,
			`CREATE TABLE IF NOT EXISTS comments (
				commentId TEXT PRIMARY KEY,
				photoId TEXT NOT NULL,
				userId TEXT NOT NULL,
				commentText TEXT NOT NULL,
				timestamp DATETIME NOT NULL,
				FOREIGN KEY(photoId) REFERENCES photos (photoId) ON DELETE CASCADE,
				FOREIGN KEY(userId) REFERENCES users (userId) ON DELETE CASCADE
			);`,
			`CREATE TABLE IF NOT EXISTS banned (
				bannerId TEXT NOT NULL,
				banneeId TEXT NOT NULL,
				PRIMARY KEY (bannerId,banneeId),
				FOREIGN KEY(bannerId) REFERENCES users (userId) ON DELETE CASCADE,
				FOREIGN KEY(banneeId) REFERENCES users (userId) ON DELETE CASCADE
			);`,
			`CREATE TABLE IF NOT EXISTS followings (
				followerId TEXT NOT NULL,
				followeeId TEXT NOT NULL,
				PRIMARY KEY (followerId,followeeId),
				FOREIGN KEY(followerId) REFERENCES users (userId) ON DELETE CASCADE,
				FOREIGN KEY(followeeId) REFERENCES users (userId) ON DELETE CASCADE
			);`,
		}

		// create them iteratively
		for i := 0; i < len(sqlStmts); i++ {
			sqlStmt := sqlStmts[i]
			_, err = db.Exec(sqlStmt)
			if err != nil {
				return nil, fmt.Errorf("error creating database structure for [%s] with error: %w", sqlStmt, err)
			}
		}
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
