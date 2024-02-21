package api

import (
	"net/http"

	"github.com/Rsen337/WASA-project-WASAPhoto-/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// banUser handles the request to ban a user.
func (rt *_router) banUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userId := ps.ByName("userId")
	otherUserId := ps.ByName("otherUserId")

	// Check if the user is authorized and authenticated to change his username.
	valid := rt.isAuthorized(userId, getToken(r.Header.Get("Authorization")))
	if valid != 0 {
		ctx.Logger.Info("banUser: User is not authorized")
		// w.WriteHeader(valid)
		writeResponse(w, valid, "User is not authorized")
		return
	}

	// Make sure that the user is not trying to ban himself.
	if userId == otherUserId {
		ctx.Logger.Info("banUser: User is trying to ban himself")
		// w.WriteHeader(http.StatusBadRequest)
		writeResponse(w, http.StatusBadRequest, "You are trying to ban yourself")
		return
	}

	err := rt.db.BanUser(userId, otherUserId)
	if err != nil {
		// The user is already banned, hence respond with 200.
		ctx.Logger.WithError(err).Error("banUser: BanUser returns error meaning the user is already banned")
		// w.WriteHeader(http.StatusOK)
		writeResponse(w, http.StatusOK, "You have already banned the user")
		return
	}

	// Remove from followings.
	err = rt.db.UnfollowUser(userId, otherUserId)
	if err != nil {
		ctx.Logger.WithError(err).Error("banUser: UnfollowUser for userId returns error")
		// w.WriteHeader(http.StatusInternalServerError)
		writeResponse(w, http.StatusInternalServerError, "")
		return
	}
	err = rt.db.UnfollowUser(otherUserId, userId)
	if err != nil {
		ctx.Logger.WithError(err).Error("banUser: UnfollowUser for otherUserId returns error")
		// w.WriteHeader(http.StatusInternalServerError)
		writeResponse(w, http.StatusInternalServerError, "")
		return
	}

	// w.WriteHeader(http.StatusCreated)
	writeResponse(w, http.StatusCreated, "You have successfully banned the user")
}
