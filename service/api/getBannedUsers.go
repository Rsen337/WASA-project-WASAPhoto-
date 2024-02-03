package api

import (
	"encoding/json"
	"github.com/Rsen337/WASA-project-WASAPhoto-/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (rt *_router) getBannedUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userId := ps.ByName("userId")
	// Parse the page query parameter from the request URL
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		// means that the page query parameter wasn't provided, hence set it to default 1
		ctx.Logger.WithError(err).Error("page query not provided")
		page = 1
	}

	// check if the user authorized and authenticated
	valid := rt.isAuthorized(userId, getToken(r.Header.Get("Authorization")))
	if valid != 0 {
		ctx.Logger.Info("uploadPhoto: isAuthorized isn't happy")
		w.WriteHeader(valid)
		return
	}

	bannedUsers, err := rt.db.GetBannedUsers(userId, page)
	if err != nil {
		ctx.Logger.WithError(err).Error("GetBannedUsers returns error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Convert photo IDs to JSON and write response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(bannedUsers); err != nil {
		ctx.Logger.WithError(err).Error("Error encoding JSON")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
