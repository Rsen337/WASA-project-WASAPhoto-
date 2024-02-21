package api

import (
	"net/http"

	"github.com/Rsen337/WASA-project-WASAPhoto-/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// followUser handles the request to follow a user.
func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userId := ps.ByName("userId")
	otherUserId := ps.ByName("otherUserId")

	// Check if the user is authorized and authenticated to change his username.
	valid := rt.isAuthorized(userId, getToken(r.Header.Get("Authorization")))
	if valid != 0 {
		ctx.Logger.Info("followUser: User is not authorized")
		// w.WriteHeader(valid)
		writeResponse(w, valid, "")
		return
	}

	// Make sure that the user is not trying to follow himself.
	if userId == otherUserId {
		ctx.Logger.Info("followUser: User is trying to follow himself")
		// w.WriteHeader(http.StatusBadRequest)
		writeResponse(w, http.StatusBadRequest, "You cannot follow yourself.")
		return
	}

	// Check if otherUserId has banned userId.
	isBanned, err := rt.db.IsBanned(otherUserId, userId)
	if err != nil {
		ctx.Logger.WithError(err).Error("followUser: IsBanned returns error")
		// w.WriteHeader(http.StatusInternalServerError)
		writeResponse(w, http.StatusInternalServerError, "")
		return
	}
	if isBanned {
		// w.WriteHeader(http.StatusForbidden)
		writeResponse(w, http.StatusForbidden, "You have been banned by this user.")
		ctx.Logger.Info("followUser: User is banned")
		return
	}

	err = rt.db.FollowUser(userId, otherUserId)
	if err != nil {
		// Check if the user is already following the other user
		if err.Error() == "UNIQUE constraint failed: followings.followerId, followings.followeeId" {
			ctx.Logger.Info("followUser: User is already following the other user")
			writeResponse(w, http.StatusOK, "You are already following the user.")
			return
		}
		ctx.Logger.WithError(err).Error("followUser: FollowUser returns error")
		// w.WriteHeader(http.StatusInternalServerError)
		writeResponse(w, http.StatusInternalServerError, "")
		return
	}

	// w.WriteHeader(http.StatusCreated)
	writeResponse(w, http.StatusCreated, "You have followed the user.")
}
