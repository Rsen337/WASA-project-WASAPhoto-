package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Rsen337/WASA-project-WASAPhoto-/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// getMyStream handles the GET request to retrieve the user's stream.
func (rt *_router) getMyStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json")

	userId := ps.ByName("userId")

	// Check if the user is authorized and authenticated.
	reqId := getToken(r.Header.Get("Authorization"))
	valid := rt.isAuthorized(userId, reqId)
	if valid != 0 {
		ctx.Logger.Info("getMyStream: User is not authorized")
		writeResponse(w, valid, "")
		return
	}

	// Parse the page query parameter from the request URL.
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		// If the page query parameter wasn't provided, set it to default 1.
		ctx.Logger.WithError(err).Error("getMyStream: Invalid page query parameter")
		page = 1
	}

	// Retrieve photos from the database.
	photos, err := rt.db.GetMyStream(userId, page)
	if err != nil {
		ctx.Logger.WithError(err).Error("getMyStream: Failed to retrieve photos from the database")
		writeResponse(w, http.StatusInternalServerError, "")
		return
	}

	// Convert photo IDs to JSON and write response.
	if err := json.NewEncoder(w).Encode(photos); err != nil {
		ctx.Logger.WithError(err).Error("getMyStream: Failed to encode JSON response")
		writeResponse(w, http.StatusInternalServerError, "Error encoding JSON")
		return
	}
}
