package api

import (
	"encoding/json"
	"net/http"

	"github.com/Rsen337/WASA-project-WASAPhoto-/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// getUserProfile handles the GET request to retrieve a user's profile.
func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userId := ps.ByName("userId")

	// Check if the user is authorized and authenticated to change their username.
	reqId := getToken(r.Header.Get("Authorization"))
	valid := rt.isAuthorized(reqId, reqId)
	if valid != 0 {
		ctx.Logger.Info("getUserProfile: User is not authorized")
		writeResponse(w, valid, "User is not authorized")
		return
	}

	// Check if userId has been banned by reqId.
	isBanned, err := rt.db.IsBanned(userId, reqId)
	if err != nil {
		ctx.Logger.WithError(err).Error("getUserProfile: IsBanned returns an error")
		writeResponse(w, http.StatusInternalServerError, "")
		return
	}
	if isBanned {
		writeResponse(w, http.StatusForbidden, "You have been banned by this user.")
		ctx.Logger.Info("getUserProfile: User is banned")
		return
	}

	user, err := rt.db.GetUserProfile(userId)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, "Error getting user profile")
		ctx.Logger.WithError(err).Error("getUserProfile: GetUserProfile returns an error")
		return
	}

	// Send the user profile to the client.
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		writeResponse(w, http.StatusInternalServerError, "Error encoding JSON")
		ctx.Logger.WithError(err).Error("getUserProfile: Error encoding JSON")
		return
	}
}
