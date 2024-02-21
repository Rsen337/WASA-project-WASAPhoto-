package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Rsen337/WASA-project-WASAPhoto-/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// getUserPhotos returns the IDs of all the photos that belong to the authenticated user.
func (rt *_router) getUserPhotos(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Get the user ID from the URL parameters
	userId := ps.ByName("userId")

	// Check if the user is authorized and authenticated
	reqId := getToken(r.Header.Get("Authorization"))
	valid := rt.isAuthorized(reqId, reqId)
	if valid != 0 {
		ctx.Logger.Info("getUserPhotos: User is not authorized")
		writeResponse(w, valid, "")
		return
	}

	// Check if the requesting user is banned
	isBanned, err := rt.db.IsBanned(userId, reqId)
	if err != nil {
		ctx.Logger.WithError(err).Error("getUserPhotos: Error checking if user is banned")
		writeResponse(w, http.StatusInternalServerError, "")
		return
	}
	if isBanned {
		writeResponse(w, http.StatusForbidden, "You have been banned by this user.")
		ctx.Logger.Info("getUserPhotos: User is banned")
		return
	}

	// Parse the page query parameter from the request URL
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		// If the page query parameter is not provided or is invalid, set it to default 1
		ctx.Logger.WithError(err).Error("getUserPhotos: Invalid or missing page query parameter")
		page = 1
	}

	var photoIds []string

	// Get the photo IDs for the user and page from the database
	photoIds, err = rt.db.GetUserPhotos(userId, page)
	if err != nil {
		ctx.Logger.WithError(err).Error("getUserPhotos: Error getting user photos from the database")
		writeResponse(w, http.StatusInternalServerError, "")
		return
	}

	// Convert photo IDs to JSON and write the response
	if err := json.NewEncoder(w).Encode(photoIds); err != nil {
		ctx.Logger.WithError(err).Error("getUserPhotos: Error encoding JSON")
		writeResponse(w, http.StatusInternalServerError, "Error encoding JSON")
		return
	}
}
