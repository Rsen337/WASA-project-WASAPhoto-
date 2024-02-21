package api

import (
	"encoding/json"
	"net/http"

	"github.com/Rsen337/WASA-project-WASAPhoto-/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// commentPhoto handles the comment photo API endpoint.
func (rt *_router) commentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	var comment Comment
	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		// Return a bad request response if the request body couldn't be decoded.
		writeResponse(w, http.StatusBadRequest, "Couldn't decode request body")
		ctx.Logger.WithError(err).Error("Couldn't decode request body")
		return
	}

	photoID := ps.ByName("photoId")

	// Check if the user is authorized and authenticated.
	reqID := getToken(r.Header.Get("Authorization"))
	valid := rt.isAuthorized(reqID, reqID)
	if valid != 0 {
		ctx.Logger.Info("commentPhoto: User is not authorized")
		// Return the authorization error.
		writeResponse(w, valid, "")
		return
	}

	err = rt.db.CommentPhoto(photoID, reqID, comment.Content)
	if err != nil {
		ctx.Logger.WithError(err).Error("CommentPhoto returned an error")
		// Return an internal server error response if CommentPhoto fails.
		writeResponse(w, http.StatusInternalServerError, "")
		return
	}

	// Return a success response.
	writeResponse(w, http.StatusCreated, "Commented on photo successfully")
}
