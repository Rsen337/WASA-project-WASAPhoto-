package api

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Rsen337/WASA-project-WASAPhoto-/service/api/reqcontext"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

// uploadPhoto handles the upload of a photo for the authenticated user.
func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get the user ID from the URL parameter
	userId := ps.ByName("userId")

	// Check if the user is authorized and authenticated to upload a photo
	valid := rt.isAuthorized(userId, getToken(r.Header.Get("Authorization")))
	if valid != 0 {
		ctx.Logger.Info("uploadPhoto: user is not authorized")
		// w.WriteHeader(valid)
		writeResponse(w, valid, "")
		return
	}

	// Read the request body
	data, err := io.ReadAll(r.Body)
	if err != nil {
		ctx.Logger.WithError(err).Error("uploadPhoto: error reading request body")
		// w.WriteHeader(http.StatusInternalServerError)
		writeResponse(w, http.StatusInternalServerError, "")
		return
	}

	// Detect the content type of the photo
	contentType := http.DetectContentType(data)

	// Check if the content type indicates a JPEG file
	if contentType == "image/jpeg" {
		ctx.Logger.Info("uploadPhoto: photo is in JPEG format")
	} else if contentType == "image/png" {
		ctx.Logger.Info("uploadPhoto: photo is in PNG format")
	} else {
		ctx.Logger.Info("uploadPhoto: photo is not in JPEG nor PNG format")
		writeResponse(w, http.StatusBadRequest, "Invalid photo format")
		return
	}

	// Replace the request body with a new buffer containing the photo data
	r.Body = io.NopCloser(bytes.NewBuffer(data))

	// Generate a unique photo ID
	photoId := uuid.New().String()

	// Create a new file for storing the photo
	out, err := os.Create(filepath.Join(photoFolder, userId, photoId))
	if err != nil {
		// w.WriteHeader(http.StatusInternalServerError)
		writeResponse(w, http.StatusInternalServerError, "")
		ctx.Logger.WithError(err).Error("uploadPhoto: error creating local photo file")
		return
	}
	defer out.Close()

	// Copy the photo data to the file
	_, err = io.Copy(out, r.Body)
	if err != nil {
		// w.WriteHeader(http.StatusInternalServerError)
		writeResponse(w, http.StatusInternalServerError, "")
		ctx.Logger.WithError(err).Error("uploadPhoto: error copying photo data to file")
		return
	}

	// Upload the photo ID and user ID to the database
	err = rt.db.UploadPhoto(photoId, userId)
	if err != nil {
		ctx.Logger.WithError(err).Error("uploadPhoto: error uploading photo to database")
		// w.WriteHeader(http.StatusInternalServerError)
		writeResponse(w, http.StatusInternalServerError, "")
		return
	}

	// Respond with success status code
	// w.WriteHeader(http.StatusCreated)
	writeResponse(w, http.StatusCreated, "")
}
