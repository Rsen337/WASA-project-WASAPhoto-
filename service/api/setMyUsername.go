package api

import (
	"encoding/json"
	"net/http"

	"github.com/Rsen337/WASA-project-WASAPhoto-/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

type Username struct {
	Username string `json:"newUsername"`
}

func (rt *_router) setMyUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	userId := ps.ByName("userId")

	// check if the user authorized and authenticated to change his username
	valid := rt.isAuthorized(userId, getToken(r.Header.Get("Authorization")))
	if valid != 0 {
		w.WriteHeader(valid)
		return
	}

	// get a new username
	var newUsername Username
	err := json.NewDecoder(r.Body).Decode(&newUsername)
	if err != nil {
		ctx.Logger.WithError(err).Error("setMyUsername: error decoding json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// check its validity
	if !usernameIsValid(newUsername.Username) {
		ctx.Logger.Infof("setMyUsername: username is not valid")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// check if the username is similar to the old one
	usernameIsSame, err := rt.db.UsernameIsSame(userId, newUsername.Username)
	if err != nil {
		ctx.Logger.WithError(err).Error("setMyUsername: UsernameIsSame returns error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if usernameIsSame {
		// Respond with 204 http status
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// check if the username is already taken
	usernameIsTaken, err := rt.db.UsernameIsTaken(userId, newUsername.Username)
	if err != nil {
		ctx.Logger.WithError(err).Error("setMyUsername: UsernameIsTaken gives error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if usernameIsTaken {
		w.WriteHeader(http.StatusConflict)
	}

	// change username in the database
	err = rt.db.SetMyUsername(userId, newUsername.Username)
	if err != nil {
		ctx.Logger.WithError(err).Error("setMyUsername: error executing update query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
