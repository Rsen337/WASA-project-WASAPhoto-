package api

import (
	"encoding/json"
	"github.com/Rsen337/WASA-project-WASAPhoto-/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) getComments(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	photoId := ps.ByName("photoId")

	// check if the user authorized and authenticated
	reqId := getToken(r.Header.Get("Authorization"))
	valid := rt.isAuthorized(reqId, reqId)
	if valid != 0 {
		ctx.Logger.Info("uploadPhoto: isAuthorized isn't happy")
		w.WriteHeader(valid)
		return
	}

	comments, err := rt.db.GetPhotoComments(photoId)
	if err != nil {
		ctx.Logger.WithError(err).Error("GetPhotoComments returns error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Convert photo IDs to JSON and write response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(comments); err != nil {
		ctx.Logger.WithError(err).Error("Error encoding JSON")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
