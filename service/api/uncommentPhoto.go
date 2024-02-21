package api

import (
	"net/http"

	"github.com/Rsen337/WASA-project-WASAPhoto-/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// uncommentPhoto handles the HTTP request to uncomment a photo.
func (rt *_router) uncommentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	photoID := ps.ByName("photoId")
	commentID := ps.ByName("commentId")

	// Check if the user is authorized and authenticated.
	reqID := getToken(r.Header.Get("Authorization"))
	isValid := rt.isAuthorized(reqID, reqID)
	if isValid != 0 {
		ctx.Logger.Info("uncommentPhoto: User is not authorized")
		writeResponse(w, isValid, "")
		return
	}

	// Check if the user is the owner of the photo.
	isOwner, err := rt.db.IsPhotoOwner(photoID, reqID)
	if err != nil {
		ctx.Logger.WithError(err).Error("uncommentPhoto: Failed to check if user is the owner of the photo")
		writeResponse(w, http.StatusInternalServerError, "")
		return
	} else if !isOwner {
		writeResponse(w, http.StatusForbidden, "You are not the owner of the photo")
		return
	}

	// Check if the comment exists.
	commentExists, err := rt.db.CommentExists(commentID)
	if err != nil {
		ctx.Logger.WithError(err).Error("uncommentPhoto: Failed to check if comment exists")
		writeResponse(w, http.StatusInternalServerError, "")
		return
	} else if !commentExists {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Uncomment the photo.
	err = rt.db.UncommentPhoto(commentID)
	if err != nil {
		ctx.Logger.WithError(err).Error("uncommentPhoto: Failed to uncomment the photo")
		writeResponse(w, http.StatusInternalServerError, "")
		return
	}

	writeResponse(w, http.StatusOK, "Uncommented on photo successfully")
}
