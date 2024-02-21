package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Rsen337/WASA-project-WASAPhoto-/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// getFollowings handles the GET request to retrieve a user's followings.
func (rt *_router) getFollowings(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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
	reqId := getToken(r.Header.Get("Authorization"))
	valid := rt.isAuthorized(reqId, reqId)
	if valid != 0 {
		ctx.Logger.Info("getFollowings: user is not authorized")
		// w.WriteHeader(valid)
		writeResponse(w, valid, "")
		return
	}

	// Check if the requesting user is banned
	isBanned, err := rt.db.IsBanned(userId, reqId)
	if err != nil {
		ctx.Logger.WithError(err).Error("IsBanned returns error")
		// w.WriteHeader(http.StatusInternalServerError)
		writeResponse(w, http.StatusInternalServerError, "")
		return
	}
	if isBanned {
		// w.WriteHeader(http.StatusForbidden)
		writeResponse(w, http.StatusForbidden, "You have been banned by this user.")
		ctx.Logger.Info("getFollowings: user is banned")
		return
	}

	// Retrieve the user's followings from the database
	followings, err := rt.db.GetFollowings(userId, page)
	if err != nil {
		ctx.Logger.WithError(err).Error("GetFollowings returns error")
		// w.WriteHeader(http.StatusInternalServerError)
		writeResponse(w, http.StatusInternalServerError, "")
		return
	}

	// Convert followings to JSON and write response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(followings); err != nil {
		ctx.Logger.WithError(err).Error("Error encoding JSON")
		// w.WriteHeader(http.StatusInternalServerError)
		writeResponse(w, http.StatusInternalServerError, "")
		return
	}
}
