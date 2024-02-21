package api

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/Rsen337/WASA-project-WASAPhoto-/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// deletePhoto handles the DELETE request to delete a photo.
func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userId := ps.ByName("userId")
	photoId := ps.ByName("photoId")

	// Check if the user is authorized and authenticated to delete the photo.
	valid := rt.isAuthorized(userId, getToken(r.Header.Get("Authorization")))
	if valid != 0 {
		ctx.Logger.Info("deletePhoto: User is not authorized")
		// w.WriteHeader(valid)
		writeResponse(w, valid, "")
		return
	}

	author, _, err := rt.db.GetPhoto(photoId)
	if err != nil {
		// The photo doesn't exist in the database.
		ctx.Logger.WithError(err).Error("deletePhoto: Photo to be deleted doesn't exist")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if author != userId {
		ctx.Logger.Info("deletePhoto: User is not authorized")
		writeResponse(w, http.StatusUnauthorized, "You are not authorized to delete this photo")
		return
	}

	// Delete the photo from the database.
	err = rt.db.DeletePhoto(photoId)
	if err != nil {
		// The photo doesn't exist in the database.
		ctx.Logger.WithError(err).Error("deletePhoto: Photo to be deleted doesn't exist")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Construct the path to the photo file.
	photoPath := filepath.Join(photoFolder, userId, photoId)

	// Delete the photo file from the filesystem.
	err = os.Remove(photoPath)
	if err != nil {
		// Failed to delete the photo file.
		ctx.Logger.WithError(err).Error("deletePhoto: Failed to delete the photo file")
		// w.WriteHeader(http.StatusInternalServerError)
		writeResponse(w, http.StatusInternalServerError, "")
		return
	}

	// Photo deleted successfully.
	// w.WriteHeader(http.StatusOK)
	writeResponse(w, http.StatusOK, "You have successfully deleted the photo")
}
