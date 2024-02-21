package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Rsen337/WASA-project-WASAPhoto-/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// getBannedUsers handles the GET request to retrieve banned users.
func (rt *_router) getBannedUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userId := ps.ByName("userId")

	// Parse the page query parameter from the request URL
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		// If the page query parameter wasn't provided, set it to default 1
		ctx.Logger.WithError(err).Error("page query not provided")
		page = 1
	}

	// Check if the user is authorized and authenticated
	valid := rt.isAuthorized(userId, getToken(r.Header.Get("Authorization")))
	if valid != 0 {
		ctx.Logger.Info("getBannedUsers: User is not authorized")
		// w.WriteHeader(valid)
		writeResponse(w, valid, "User is not authorized")
		return
	}

	// Retrieve banned users from the database
	bannedUsers, err := rt.db.GetBannedUsers(userId, page)
	if err != nil {
		ctx.Logger.WithError(err).Error("GetBannedUsers returns error")
		// w.WriteHeader(http.StatusInternalServerError)
		writeResponse(w, http.StatusInternalServerError, "")
		return
	}

	// Set the response header to indicate JSON content
	w.Header().Set("Content-Type", "application/json")

	// Encode the banned users to JSON and write the response
	if err := json.NewEncoder(w).Encode(bannedUsers); err != nil {
		ctx.Logger.WithError(err).Error("Error encoding JSON")
		// w.WriteHeader(http.StatusInternalServerError)
		writeResponse(w, http.StatusInternalServerError, "Error encoding JSON")
		return
	}
}
