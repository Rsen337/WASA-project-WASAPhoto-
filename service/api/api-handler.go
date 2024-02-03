package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here

func (rt *_router) Handler() http.Handler {

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	// Register or log in
	rt.router.POST("/session", rt.wrap(rt.doLogin))

	// Change username
	rt.router.PUT("/users/:userId/change-username", rt.wrap(rt.setMyUsername))

	// Upload or delete my photos
	rt.router.GET("/users/:userId/photo", rt.wrap(rt.getUserPhotos))
	rt.router.POST("/users/:userId/photo", rt.wrap(rt.uploadPhoto))
	rt.router.DELETE("/users/:userId/photo/:photoId", rt.wrap(rt.deletePhoto))

	// manage followers
	rt.router.GET("/users/:userId/followers", rt.wrap(rt.getFollowers))

	// manage followings
	rt.router.GET("/users/:userId/followings", rt.wrap(rt.getFollowings))
	rt.router.PUT("/users/:userId/followings/:otherUserId", rt.wrap(rt.followUser))
	rt.router.DELETE("/users/:userId/followings/:otherUserId", rt.wrap(rt.unfollowUser))

	// manage bans
	rt.router.GET("/users/:userId/banned", rt.wrap(rt.getBannedUsers))
	rt.router.PUT("/users/:userId/banned/:otherUserId", rt.wrap(rt.banUser))
	rt.router.DELETE("/users/:userId/banned/:otherUserId", rt.wrap(rt.unbanUser))

	// search users
	rt.router.GET("/users", rt.wrap(rt.searchUser))
	rt.router.GET("/users/:userId", rt.wrap(rt.getUserProfile))

	// stream of photos
	rt.router.GET("/users/:userId/stream", rt.wrap(rt.getMyStream))

	// photos
	rt.router.GET("/photos/:photoId", rt.wrap(rt.getPhoto))
	rt.router.GET("/photos/:photoId/details", rt.wrap(rt.getPhotoDetails))

	// likes
	rt.router.PUT("/photos/:photoId/like/:userId", rt.wrap(rt.likePhoto))
	rt.router.DELETE("/photos/:photoId/like/:userId", rt.wrap(rt.unlikePhoto))

	// comments
	rt.router.GET("/photos/:photoId/comment", rt.wrap(rt.getComments))
	rt.router.POST("/photos/:photoId/comment", rt.wrap(rt.commentPhoto))
	rt.router.DELETE("/photos/:photoId/comment/:commentId", rt.wrap(rt.uncommentPhoto))

	return rt.router
}
