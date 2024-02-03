package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Rsen337/WASA-project-WASAPhoto-/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getFollowers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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
	reqId := getToken(r.Header.Get("Authorization"))
	valid := rt.isAuthorized(reqId, reqId)
	if valid != 0 {
		ctx.Logger.Info("uploadPhoto: isAuthorized isn't happy")
		w.WriteHeader(valid)
		return
	}

	// checks if the requesting user is banned
	isBanned, err := rt.db.IsBanned(userId, reqId)
	if err != nil {
		ctx.Logger.WithError(err).Error("IsBanned returns error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if isBanned {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.Info("is banned")
		return
	}

	followings, err := rt.db.GetFollowers(userId, page)
	if err != nil {
		ctx.Logger.WithError(err).Error("GetFollowers returns error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Convert photo IDs to JSON and write response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(followings); err != nil {
		ctx.Logger.WithError(err).Error("Error encoding JSON")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
