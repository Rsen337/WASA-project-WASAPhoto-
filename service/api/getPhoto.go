package api

import (
	"encoding/json"
	"net/http"
	"path/filepath"

	"github.com/Rsen337/WASA-project-WASAPhoto-/service/api/reqcontext"
	"github.com/Rsen337/WASA-project-WASAPhoto-/service/database"
	"github.com/julienschmidt/httprouter"
)

// get a specific photo searched up by the photoId
func (rt *_router) getPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	valid := rt.isAuthorized(getToken(r.Header.Get("Authorization")), getToken(r.Header.Get("Authorization")))
	if valid != 0 {
		ctx.Logger.Info("uploadPhoto: isAuthorized isn't happy")
		w.WriteHeader(valid)
		return
	}

	photoId := ps.ByName("photoId")

	userId, _, err := rt.db.GetPhoto(photoId)
	if err != nil {
		// the requested photo doesn't exist
		w.WriteHeader(http.StatusNoContent)
		ctx.Logger.WithError(err).Error("GetPhotoDetails retuns an error")
		return
	}

	http.ServeFile(w, r,
		filepath.Join(photoFolder, userId, photoId))
}

// get details of a specific photo searched up by the photoId
func (rt *_router) getPhotoDetails(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// check if the user authorized and authenticated to change his username
	valid := rt.isAuthorized(getToken(r.Header.Get("Authorization")), getToken(r.Header.Get("Authorization")))
	if valid != 0 {
		ctx.Logger.Info("uploadPhoto: isAuthorized isn't happy")
		w.WriteHeader(valid)
		return
	}

	photoId := ps.ByName("photoId")

	userId, timestamp, err := rt.db.GetPhoto(photoId)
	if err != nil {
		// the requested photo doesn't exist
		w.WriteHeader(http.StatusNoContent)
		ctx.Logger.WithError(err).Error("GetPhotoDetails retuns an error")
		return
	}

	likesAmount, err := rt.db.GetPhotoLikes(photoId)

	comments, err := rt.db.GetPhotoComments(photoId)

	w.Header().Set("Content-Type", "application/json")

	photo := database.Photo{
		PhotoID:     photoId,
		UserID:      userId,
		Timestamp:   timestamp,
		LikesAmount: likesAmount,
		Comments:    comments,
	}
	err = json.NewEncoder(w).Encode(photo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("session: can't create response json")
		return
	}
}
