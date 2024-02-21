package api

import (
	"encoding/json"
	"net/http"

	"github.com/Rsen337/WASA-project-WASAPhoto-/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// Username represents the new username to be set
type Username struct {
	Username string `json:"newUsername"`
}

// setMyUsername is the handler function for setting the username
func (rt *_router) setMyUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userId := ps.ByName("userId")

	// Check if the user is authorized and authenticated to change their username
	valid := rt.isAuthorized(userId, getToken(r.Header.Get("Authorization")))
	if valid != 0 {
		// w.WriteHeader(valid)
		writeResponse(w, valid, "")
		return
	}

	// Decode the new username from the request body
	var newUsername Username
	err := json.NewDecoder(r.Body).Decode(&newUsername)
	if err != nil {
		ctx.Logger.WithError(err).Error("setMyUsername: error decoding JSON")
		// w.WriteHeader(http.StatusBadRequest)
		writeResponse(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	// Check the validity of the new username
	if !usernameIsValid(newUsername.Username) {
		ctx.Logger.Infof("setMyUsername: username is not valid")
		// w.WriteHeader(http.StatusBadRequest)
		writeResponse(w, http.StatusBadRequest, "Invalid username")
		return
	}

	// Check if the new username is the same as the old one
	usernameIsSame, err := rt.db.UsernameIsSame(userId, newUsername.Username)
	if err != nil {
		ctx.Logger.WithError(err).Error("setMyUsername: UsernameIsSame returns an error")
		// w.WriteHeader(http.StatusInternalServerError)
		writeResponse(w, http.StatusInternalServerError, "")
		return
	} else if usernameIsSame {
		// Respond with 204 HTTP status (No Content)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Check if the new username is already taken
	usernameIsTaken, err := rt.db.UsernameIsTaken(userId, newUsername.Username)
	if err != nil {
		ctx.Logger.WithError(err).Error("setMyUsername: UsernameIsTaken returns an error")
		// w.WriteHeader(http.StatusInternalServerError)
		writeResponse(w, http.StatusInternalServerError, "")
		return
	} else if usernameIsTaken {
		// w.WriteHeader(http.StatusConflict)
		writeResponse(w, http.StatusConflict, "Username is already taken")
		return
	}

	// Change the username in the database
	err = rt.db.SetMyUsername(userId, newUsername.Username)
	if err != nil {
		ctx.Logger.WithError(err).Error("setMyUsername: error executing update query")
		// w.WriteHeader(http.StatusInternalServerError)
		writeResponse(w, http.StatusInternalServerError, "")
		return
	}

	// w.WriteHeader(http.StatusOK)
	writeResponse(w, http.StatusOK, "Username updated successfully")

}
