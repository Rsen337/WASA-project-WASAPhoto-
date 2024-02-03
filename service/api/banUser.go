package api

import (
	"net/http"

	"github.com/Rsen337/WASA-project-WASAPhoto-/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) banUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userId := ps.ByName("userId")
	otherUserId := ps.ByName("otherUserId")

	// check if the user authorized and authenticated to change his username
	valid := rt.isAuthorized(userId, getToken(r.Header.Get("Authorization")))
	if valid != 0 {
		ctx.Logger.Info("uploadPhoto: isAuthorized isn't happy")
		w.WriteHeader(valid)
		return
	}

	// make sure that the user is not trying to ban himself
	if userId == otherUserId {
		ctx.Logger.Info("user is trying to ban himself")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := rt.db.BanUser(userId, otherUserId)
	if err != nil {
		// already banned, hence respond with 200
		ctx.Logger.WithError(err).Error("BanUser returns error")
		w.WriteHeader(http.StatusOK)
		return
	}

	// remove from followings
	err = rt.db.UnfollowUser(userId, otherUserId)
	if err != nil {
		ctx.Logger.WithError(err).Error("UnfollowUser for userId returns error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = rt.db.UnfollowUser(otherUserId, userId)
	if err != nil {
		ctx.Logger.WithError(err).Error("UnfollowUser for otherUserId returns error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
