package api

import (
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// path to the folder that contains all the photos
var photoFolder = filepath.Join("/tmp", "media")

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
	Content string `json:"commentText"`
}

// Checks if the username is valid
func usernameIsValid(username string) bool {
	// Define the regular expression pattern for the username
	pattern := "^[a-zA-Z0-9_-]{3,32}$"

	// Compile the regular expression
	regex, err := regexp.Compile(pattern)
	if err != nil {
		// Handle error if the pattern is invalid
		return false
	}

	// Check if the username matches the pattern
	return regex.MatchString(username)
}

func (rt *_router) isAuthorized(reqUserId string, auth string) int {

	if auth == "" {
		return http.StatusForbidden
	}

	reqExists, err := rt.db.UserExists(reqUserId)
	if err != nil {
		return http.StatusInternalServerError
	} else if !reqExists {
		return http.StatusForbidden
	}

	authExists, err := rt.db.UserExists(auth)
	if err != nil {
		return http.StatusInternalServerError
	} else if !authExists {
		return http.StatusForbidden
	}

	if reqUserId != auth {
		return http.StatusUnauthorized
	}

	return 0
}

// extract token from the authorization header
func getToken(authorization string) string {
	var tokens = strings.Split(authorization, " ")
	if len(tokens) == 2 {
		return strings.Trim(tokens[1], " ")
	}
	return ""
}
