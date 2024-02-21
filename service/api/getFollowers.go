package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Rsen337/WASA-project-WASAPhoto-/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// getFollowers handles the GET request to retrieve followers of a user.
func (rt *_router) getFollowers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userId := ps.ByName("userId")

	// Parse the page query parameter from the request URL
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		// If the page query parameter wasn't provided or is invalid, set it to default 1
		ctx.Logger.WithError(err).Error("Invalid or missing page query parameter")
		page = 1
	}

	// Check if the user is authorized and authenticated
	reqId := getToken(r.Header.Get("Authorization"))
	valid := rt.isAuthorized(reqId, reqId)
	if valid != 0 {
		ctx.Logger.Info("getFollowers: User is not authorized")
		// w.WriteHeader(valid)
		writeResponse(w, valid, "")
		return
	}

	// Check if the requesting user is banned
	isBanned, err := rt.db.IsBanned(userId, reqId)
	if err != nil {
		ctx.Logger.WithError(err).Error("IsBanned returns an error")
		// w.WriteHeader(http.StatusInternalServerError)
		writeResponse(w, http.StatusInternalServerError, "")
		return
	}
	if isBanned {
		// w.WriteHeader(http.StatusForbidden)
		writeResponse(w, http.StatusForbidden, "You have been banned by this user.")
		ctx.Logger.Info("getFollowers: User is banned")
		return
	}

	// Retrieve followers from the database
	followers, err := rt.db.GetFollowers(userId, page)
	if err != nil {
		ctx.Logger.WithError(err).Error("GetFollowers returns an error")
		// w.WriteHeader(http.StatusInternalServerError)
		writeResponse(w, http.StatusInternalServerError, "")
		return
	}

	// Convert followers to JSON and write the response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(followers); err != nil {
		ctx.Logger.WithError(err).Error("Error encoding JSON")
		// w.WriteHeader(http.StatusInternalServerError)
		writeResponse(w, http.StatusInternalServerError, "")
		return
	}
}
