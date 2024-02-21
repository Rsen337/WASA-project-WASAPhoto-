package api

import (
	"net/http"

	"github.com/Rsen337/WASA-project-WASAPhoto-/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// unfollowUser handles the HTTP request to unfollow a user.
func (rt *_router) unfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userId := ps.ByName("userId")
	otherUserId := ps.ByName("otherUserId")

	// Check if the user is authorized and authenticated to unfollow the other user.
	valid := rt.isAuthorized(userId, getToken(r.Header.Get("Authorization")))
	if valid != 0 {
		ctx.Logger.Info("unfollowUser: User is not authorized")
		// w.WriteHeader(valid)
		writeResponse(w, valid, "User is not authorized")
		return
	}

	err := rt.db.UnfollowUser(userId, otherUserId)
	if err != nil {
		ctx.Logger.WithError(err).Error("unfollowUser: Failed to unfollow user")
		// w.WriteHeader(http.StatusInternalServerError)
		writeResponse(w, http.StatusInternalServerError, "Failed to unfollow user.")
		return
	}

	// w.WriteHeader(http.StatusOK)
	writeResponse(w, http.StatusOK, "You have unfollowed the user.")
}
