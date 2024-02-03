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

// returns ids of all the photos that belong to the athenticated user
func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// w.Header().Set("Content-Type", "application/json")

	userId := ps.ByName("userId")

	// check if the user authorized and authenticated to change his username
	valid := rt.isAuthorized(userId, getToken(r.Header.Get("Authorization")))
	if valid != 0 {
		ctx.Logger.Info("uploadPhoto: isAuthorized isn't happy")
		w.WriteHeader(valid)
		return
	}

	// Create a copy of the body
	data, err := io.ReadAll(r.Body)
	if err != nil {
		ctx.Logger.WithError(err).Error("uploadPhoto: error reading body content")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Detect the content type
	contentType := http.DetectContentType(data)

	// Check if the content type indicates a JPEG file
	if contentType == "image/jpeg" {
		ctx.Logger.Info("it is jpeg")
	} else if contentType == "image/png" {
		ctx.Logger.Info("it is png")
	} else {
		ctx.Logger.Info("it is not jpeg nor png")
	}

	r.Body = io.NopCloser(bytes.NewBuffer(data))

	// Generate a unique photo ID
	photoId := uuid.New().String()

	// Create an empty file for storing the body content (image)
	out, err := os.Create(filepath.Join(photoFolder, userId, photoId))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("uploadPhoto: error creating local photo file")
		return
	}

	// Copy body content to the previously created file
	_, err = io.Copy(out, r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("uploadPhoto: error copying body content into file photo")
		return
	}

	// Close the created file
	out.Close()

	err = rt.db.UploadPhoto(photoId, userId)
	if err != nil {
		ctx.Logger.WithError(err).Error("rt.db.UploadPhoto returns error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Respond with success message
	w.WriteHeader(http.StatusCreated)

}
