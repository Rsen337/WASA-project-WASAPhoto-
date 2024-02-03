package api

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/Rsen337/WASA-project-WASAPhoto-/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	userId := ps.ByName("userId")
	photoId := ps.ByName("photoId")

	// check if the user authorized and authenticated to change his username
	valid := rt.isAuthorized(userId, getToken(r.Header.Get("Authorization")))
	if valid != 0 {
		ctx.Logger.Info("uploadPhoto: isAuthorized isn't happy")
		w.WriteHeader(valid)
		return
	}

	// deletes the photo from the database and returns and error if it doesn't exist
	err := rt.db.DeletePhoto(photoId)
	if err != nil {
		// the photo didn't exist from the start
		ctx.Logger.WithError(err).Error("photo to be deleted doesn't even exist")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Construct the path to the photo file
	photoPath := filepath.Join(photoFolder, userId, photoId)

	// Delete the photo file from the filesystem
	err = os.Remove(photoPath)
	if err != nil {
		// Failed to delete the photo file
		ctx.Logger.WithError(err).Error("photo doesn't exist")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Photo deleted successfully
	w.WriteHeader(http.StatusOK)

}
