package api

import (
	"github.com/Rsen337/WASA-project-WASAPhoto-/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) unlikePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	photoId := ps.ByName("photoId")
	userId := ps.ByName("userId")

	// check if the user authorized and authenticated
	reqId := getToken(r.Header.Get("Authorization"))
	valid := rt.isAuthorized(userId, reqId)
	if valid != 0 {
		ctx.Logger.Info("uploadPhoto: isAuthorized isn't happy")
		w.WriteHeader(valid)
		return
	}

	err := rt.db.UnlikePhoto(photoId, userId)
	if err != nil {
		ctx.Logger.WithError(err).Error("UnlikePhoto returns error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
