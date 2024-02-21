package api

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Rsen337/WASA-project-WASAPhoto-/service/api/reqcontext"
	"github.com/Rsen337/WASA-project-WASAPhoto-/service/database"
	"github.com/julienschmidt/httprouter"
)

// getPhoto retrieves a specific photo searched by the photoId.
func (rt *_router) getPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Check if the request is authorized.
	valid := rt.isAuthorized(getToken(r.Header.Get("Authorization")), getToken(r.Header.Get("Authorization")))
	if valid != 0 {
		ctx.Logger.Info("getPhoto: Unauthorized request")
		writeResponse(w, valid, "")
		return
	}

	photoId := ps.ByName("photoId")

	userId, _, err := rt.db.GetPhoto(photoId)
	if err != nil {
		// The requested photo doesn't exist.
		w.WriteHeader(http.StatusNoContent)
		ctx.Logger.WithError(err).Error("getPhoto: GetPhotoDetails returns an error")
		return
	}

	photoPath := filepath.Join(photoFolder, userId, photoId)

	// Open the photo file
	file, err := os.Open(photoPath)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, "")
		ctx.Logger.WithError(err).Error("getPhoto: Failed to open photo file")
		return
	}
	defer file.Close()

	// Get the file information
	fileInfo, err := file.Stat()
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, "")
		ctx.Logger.WithError(err).Error("getPhoto: Failed to get file information")
		return
	}

	// Set the response headers
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))

	// Serve the file content
	http.ServeContent(w, r, photoId, fileInfo.ModTime(), file)
}

// getPhotoDetails retrieves details of a specific photo searched by the photoId.
func (rt *_router) getPhotoDetails(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Check if the request is authorized.
	valid := rt.isAuthorized(getToken(r.Header.Get("Authorization")), getToken(r.Header.Get("Authorization")))
	if valid != 0 {
		ctx.Logger.Info("getPhotoDetails: Unauthorized request")
		writeResponse(w, valid, "")
		return
	}

	photoId := ps.ByName("photoId")

	userId, timestamp, err := rt.db.GetPhoto(photoId)
	if err != nil {
		// The requested photo doesn't exist.
		w.WriteHeader(http.StatusNoContent)
		ctx.Logger.WithError(err).Error("getPhotoDetails: GetPhotoDetails returns an error")
		return
	}

	likesAmount, err := rt.db.GetPhotoLikes(photoId)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, "")
		ctx.Logger.WithError(err).Error("getPhotoDetails: GetPhotoLikes returns an error")
		return
	}

	comments, err := rt.db.GetPhotoComments(photoId, 1)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, "")
		ctx.Logger.WithError(err).Error("getPhotoDetails: GetPhotoComments returns an error")
		return
	}

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
		writeResponse(w, http.StatusInternalServerError, "Error encoding JSON")
		ctx.Logger.WithError(err).Error("getPhotoDetails: Can't create response JSON")
		return
	}
}
