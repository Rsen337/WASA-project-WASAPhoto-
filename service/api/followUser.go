package api

import (
	"net/http"

	"github.com/Rsen337/WASA-project-WASAPhoto-/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	userId := ps.ByName("userId")
	otherUserId := ps.ByName("otherUserId")

	// check if the user authorized and authenticated to change his username
	valid := rt.isAuthorized(userId, getToken(r.Header.Get("Authorization")))
	if valid != 0 {
		ctx.Logger.Info("uploadPhoto: isAuthorized isn't happy")
		w.WriteHeader(valid)
		return
	}

	// make sure that the user is not trying to follow himself
	if userId == otherUserId {
		ctx.Logger.Info("user is trying to follow himself")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// check if otherUserId has banned userId
	isBanned, err := rt.db.IsBanned(otherUserId, userId)
	if err != nil {
		ctx.Logger.WithError(err).Error("IsBanned returns error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if isBanned {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.Info("is banned")
		return
	}

	err = rt.db.FollowUser(userId, otherUserId)
	if err != nil {
		ctx.Logger.WithError(err).Error("FollowUser returns error")
		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
