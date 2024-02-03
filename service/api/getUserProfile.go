package api

import (
	"encoding/json"
	"github.com/Rsen337/WASA-project-WASAPhoto-/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	userId := ps.ByName("userId")

	// check if the user authorized and authenticated to change his username
	valid := rt.isAuthorized(getToken(r.Header.Get("Authorization")), getToken(r.Header.Get("Authorization")))
	if valid != 0 {
		ctx.Logger.Info("uploadPhoto: isAuthorized isn't happy")
		w.WriteHeader(valid)
		return
	}

	user, err := rt.db.GetUserProfile(userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("GetUserProfile returns an error")
		return
	}

	// send the list to the client
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		ctx.Logger.WithError(err).Error("Error encoding JSON")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
