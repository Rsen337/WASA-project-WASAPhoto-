package api

import (
	"github.com/Rsen337/WASA-project-WASAPhoto-/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// likePhoto handles the HTTP POST request for liking a photo.
func (rt *_router) likePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	photoID := ps.ByName("photoId")
	userID := ps.ByName("userId")

	// Check if the user is authorized and authenticated.
	reqID := getToken(r.Header.Get("Authorization"))
	isAuthorized := rt.isAuthorized(userID, reqID)
	if isAuthorized != 0 {
		ctx.Logger.Info("likePhoto: User is not authorized")
		writeResponse(w, isAuthorized, "")
		return
	}

	err := rt.db.LikePhoto(photoID, userID)
	if err != nil {
		// The photo is already liked.
		ctx.Logger.WithError(err).Error("likePhoto: LikePhoto returns an error")
		writeResponse(w, http.StatusOK, "You already liked this photo")
		return
	}

	writeResponse(w, http.StatusCreated, "Photo liked successfully")
}
