package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Rsen337/WASA-project-WASAPhoto-/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getComments(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	photoID := ps.ByName("photoId")

	// Parse the page query parameter from the request URL
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		// If the page query parameter wasn't provided or is invalid, set it to default 1
		ctx.Logger.WithError(err).Error("Invalid or missing page query parameter")
		page = 1
	}

	// Check if the user is authorized and authenticated
	reqID := getToken(r.Header.Get("Authorization"))
	valid := rt.isAuthorized(reqID, reqID)
	if valid != 0 {
		ctx.Logger.Info("getComments: User is not authorized")
		writeResponse(w, valid, "")
		return
	}

	comments, err := rt.db.GetPhotoComments(photoID, page)
	if err != nil {
		ctx.Logger.WithError(err).Error("GetPhotoComments returns error")
		writeResponse(w, http.StatusInternalServerError, "")
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")

	// Encode comments to JSON and write response
	if err := json.NewEncoder(w).Encode(comments); err != nil {
		ctx.Logger.WithError(err).Error("Error encoding JSON")
		writeResponse(w, http.StatusInternalServerError, "Error encoding JSON")
		return
	}
}
