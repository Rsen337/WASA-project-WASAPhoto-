package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Rsen337/WASA-project-WASAPhoto-/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

type Token struct {
	Token string `json:"token"` // User's token
}

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("content-type", "application/json")

	var user Token
	var token string

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("couldn't decode body")
		return
	} else if token = user.Token; !usernameIsValid(token) {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.Infof("usernameIsValid is saying nuh-uh")
		return
	}

	userId, err := rt.db.DoLogin(token)
	if err != nil {
		// something went wrong, I dunno what
		ctx.Logger.WithError(err).Error("rt.db.DoLogin returned an error")
		return
	}

	// Send the output to the client
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("session: can't create response json")
		return
	}
}

// Checks if the username
func usernameIsValid(username string) bool {
	var trimmed = strings.TrimSpace(username)
	return len(username) >= 3 && len(username) <= 32 && trimmed != "" && !strings.ContainsAny(trimmed, "?_")
}
