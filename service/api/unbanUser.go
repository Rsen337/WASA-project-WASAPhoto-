package api

import (
	"net/http"

	"github.com/Rsen337/WASA-project-WASAPhoto-/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// unbanUser handles the HTTP request to unban a user.
func (rt *_router) unbanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userId := ps.ByName("userId")
	otherUserId := ps.ByName("otherUserId")

	// Check if the user is authorized and authenticated to unban the user.
	valid := rt.isAuthorized(userId, getToken(r.Header.Get("Authorization")))
	if valid != 0 {
		ctx.Logger.Info("unbanUser: User is not authorized")
		writeResponse(w, valid, "")
		return
	}

	err := rt.db.UnbanUser(userId, otherUserId)
	if err != nil {
		ctx.Logger.WithError(err).Error("unbanUser: Failed to unban user")
		writeResponse(w, http.StatusInternalServerError, "Failed to unban user.")
		return
	}

	writeResponse(w, http.StatusOK, "You have successfully unbanned the user.")
}
