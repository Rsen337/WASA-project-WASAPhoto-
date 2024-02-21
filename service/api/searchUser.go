package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Rsen337/WASA-project-WASAPhoto-/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// searchUser handles the search user API endpoint.
func (rt *_router) searchUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Parse the "username" query parameter from the request URL
	toSearch := r.URL.Query().Get("username")

	// Parse the "page" query parameter from the request URL
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		// If the "page" query parameter wasn't provided or is invalid, set it to default 1
		ctx.Logger.WithError(err).Error("page query not provided")
		page = 1
	}

	// Check if the user is authorized and authenticated to change their username
	userId := getToken(r.Header.Get("Authorization"))
	valid := rt.isAuthorized(userId, userId)
	if valid != 0 {
		ctx.Logger.Info("searchUser: isAuthorized isn't happy")
		writeResponse(w, valid, "User is not authorized")
		return
	}

	// Get matching results
	users, err := rt.db.SearchUser(toSearch, userId, page)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, "Error searching for users")
		ctx.Logger.WithError(err).Error("SearchUser returns an error")
		return
	}
	if len(users) == 0 {
		w.WriteHeader(http.StatusNoContent)
		ctx.Logger.Info("No matching users")
		return
	}

	// Send the list of users to the client
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		ctx.Logger.WithError(err).Error("Error encoding JSON")
		writeResponse(w, http.StatusInternalServerError, "Error encoding JSON")
		return
	}
}
