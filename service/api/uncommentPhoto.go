package api

import (
	"net/http"

	"github.com/Rsen337/WASA-project-WASAPhoto-/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) uncommentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	photoId := ps.ByName("photoId")
	commentId := ps.ByName("commentId")

	// check if the user authorized and authenticated
	reqId := getToken(r.Header.Get("Authorization"))
	valid := rt.isAuthorized(reqId, reqId)
	if valid != 0 {
		ctx.Logger.Info("uploadPhoto: isAuthorized isn't happy")
		w.WriteHeader(valid)
		return
	}

	isOwner, err := rt.db.IsPhotoOwner(photoId, reqId)
	if err != nil {
		ctx.Logger.WithError(err).Error("IsPhotoOwner returns error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if !isOwner {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	commentExists, err := rt.db.CommentExists(commentId)
	if err != nil {
		ctx.Logger.WithError(err).Error("IsPhotoOwner returns error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if !commentExists {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = rt.db.UncommentPhoto(commentId)
	if err != nil {
		ctx.Logger.WithError(err).Error("CommentPhoto returns error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
