package api

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

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

	userId, exist, err := rt.db.DoLogin(token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("rt.db.DoLogin returned an error")
		return
	}
	if exist {
		w.WriteHeader(http.StatusOK)
		response := map[string]string{"identifier": userId}
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("session: can't create response json")
			return
		}
		return
	}

	// create photo folder for the user
	// Define file path
	filePath := filepath.Join(photoFolder, userId)
	err = os.MkdirAll(filePath, os.ModePerm)
	if err != nil {
		ctx.Logger.WithError(err).Error("doLogin: Error creating photo directory")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Send the output to the client
	w.WriteHeader(http.StatusCreated)
	response := map[string]string{"identifier": userId}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("session: can't create response json")
		return
	}
}
