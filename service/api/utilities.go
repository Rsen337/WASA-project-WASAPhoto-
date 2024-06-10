package api

import (
	"log"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
)

// photoFolder represents the path to the folder that contains all the photos.
var photoFolder = filepath.Join("/tmp", "media")

// Comment represents the structure of a comment.
type Comment struct {
	Content string `json:"commentText"`
}

// usernameIsValid checks if the username is valid.
// It returns true if the username matches the regular expression pattern,
// otherwise it returns false.
func usernameIsValid(username string) bool {
	// Define the regular expression pattern for the username.
	pattern := "^[a-zA-Z0-9_-]{3,32}$"

	// Compile the regular expression.
	regex, err := regexp.Compile(pattern)
	if err != nil {
		// Handle error if the pattern is invalid.
		return false
	}

	// Check if the username matches the pattern.
	return regex.MatchString(username)
}

// isAuthorized checks if the request is authorized.
// It takes the request user ID and the authorization token as parameters.
// It returns an HTTP status code indicating the authorization result.
func (rt *_router) isAuthorized(reqUserId string, auth string) int {
	if auth == "" {
		return http.StatusUnauthorized
	}

	reqExists, err := rt.db.UserExists(reqUserId)
	if err != nil {
		return http.StatusInternalServerError
	} else if !reqExists {
		return http.StatusUnauthorized
	}

	authExists, err := rt.db.UserExists(auth)
	if err != nil {
		return http.StatusInternalServerError
	} else if !authExists {
		return http.StatusUnauthorized
	}

	if reqUserId != auth {
		return http.StatusForbidden
	}

	return 0
}

// getToken extracts the token from the authorization header.
// It takes the authorization header as a parameter.
// It returns the extracted token.
func getToken(authorization string) string {
	var tokens = strings.Split(authorization, " ")
	if len(tokens) == 2 {
		return strings.Trim(tokens[1], " ")
	}
	return ""
}

func writeResponse(w http.ResponseWriter, status int, message string) {

	if message == "" {
		switch status {
		case http.StatusOK:
			message = "OK"
		case http.StatusCreated:
			message = "Created"
		case http.StatusNoContent:
			message = "No Content"
		case http.StatusBadRequest:
			message = "Bad Request"
		case http.StatusUnauthorized:
			message = "Unauthorized"
		case http.StatusForbidden:
			message = "Forbidden"
		case http.StatusNotFound:
			message = "Not Found"
		case http.StatusInternalServerError:
			message = "Internal Server Error"
		default:
			message = "Unknown Status"
		}
	}

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "text/plain")

	_, err := w.Write([]byte(message))
	if err != nil {
		log.Println(err)
	}

}
