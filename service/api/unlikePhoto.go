package api

import (
	"github.com/Rsen337/WASA-project-WASAPhoto-/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// unlikePhoto handles the unlike photo request.
func (rt *_router) unlikePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	photoID := ps.ByName("photoId")
	userID := ps.ByName("userId")

	// Check if the user is authorized and authenticated.
	reqID := getToken(r.Header.Get("Authorization"))
	isValid := rt.isAuthorized(userID, reqID)
	if isValid != 0 {
		ctx.Logger.Info("unlikePhoto: User is not authorized")
		writeResponse(w, isValid, "")
		return
	}

	err := rt.db.UnlikePhoto(photoID, userID)
	if err != nil {
		ctx.Logger.WithError(err).Error("unlikePhoto: Failed to unlike photo")
		writeResponse(w, http.StatusInternalServerError, "")
		return
	}

	writeResponse(w, http.StatusOK, "Photo unliked successfully")
}
