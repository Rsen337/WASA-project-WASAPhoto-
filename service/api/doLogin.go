package api

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Rsen337/WASA-project-WASAPhoto-/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// Token represents the user's token
type Token struct {
	Token string `json:"token"` // User's token
}

// doLogin is the handler function for the login endpoint
func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Set the content type header to JSON
	w.Header().Set("Content-Type", "application/json")

	var user Token
	var token string

	// Decode the request body into the user struct
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		// Return a bad request response with an error message
		writeResponse(w, http.StatusBadRequest, "Invalid request body")
		ctx.Logger.WithError(err).Error("doLogin: couldn't decode request body")
		return
	} else if token = user.Token; !usernameIsValid(token) {
		// Return a bad request response with an error message
		writeResponse(w, http.StatusBadRequest, "Invalid username")
		ctx.Logger.Infof("doLogin: usernameIsValid returned false")
		return
	}

	// Perform login operation in the database
	userId, exist, err := rt.db.DoLogin(token)
	if err != nil {
		// Return an internal server error response with an error message
		writeResponse(w, http.StatusInternalServerError, "Internal Server Error")
		ctx.Logger.WithError(err).Error("doLogin: rt.db.DoLogin returned an error")
		return
	}

	// If user exists, return a success response with the user identifier
	if exist {
		w.WriteHeader(http.StatusOK)
		response := map[string]string{"identifier": userId}
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			// Return an internal server error response with an error message
			writeResponse(w, http.StatusInternalServerError, "Internal Server Error")
			ctx.Logger.WithError(err).Error("doLogin: couldn't create response JSON")
			return
		}
		return
	}

	// Create photo folder for the user
	// Define the file path
	filePath := filepath.Join(photoFolder, userId)
	err = os.MkdirAll(filePath, os.ModePerm)
	if err != nil {
		ctx.Logger.WithError(err).Error("doLogin: error creating photo directory")
		// Return an internal server error response with an error message
		writeResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// Send the output to the client
	w.WriteHeader(http.StatusCreated)
	response := map[string]string{"identifier": userId}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		// Return an internal server error response with an error message
		writeResponse(w, http.StatusInternalServerError, "Internal Server Error")
		ctx.Logger.WithError(err).Error("doLogin: couldn't create response JSON")
		return
	}
}
