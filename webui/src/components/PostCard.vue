<script>
export default {
	props: ["photoID"],
	data: function () {
		return {
			// For comments
			page: 1,

			// Whether the image is loaded (disables the spinner)
			imageReady: false,

			// likesAmount and comments and other data
			post_liked: false,
			post_like_cnt: 0,
			post_comments_cnt: 0,
			comments_data: [],
			comments_shown: false,
			commentMsg: "",
			timestamp: "",
			username: "",
			author: "",
			isMyPhoto: false,

			// Whether the comments have ended (no more comments to load)
			data_ended: false,
		}
	},
	methods: {
		// Visit the user's profile
		visitUser() {
			this.$router.push({ path: "/profile/" + this.author });
		},

		// Post a new comment
		postComment() {
			this.$axios.post("/photos/" + this.photoID + "/comment", {
				"commentText": this.commentMsg,
			}).then(response => {
				if (response == null) return // the interceptors returns null if something goes bad

				// Reset the comment input and uptimestamp the counter
				this.commentMsg = "";
				this.post_comments_cnt++;

				// Fetch comments from the server
				this.comments_data = [];
				this.page = 1;
				this.getComments();
			})
		},

		// Delete a comment
		deleteComment(commentID) {
			this.$axios.delete("/photos/" + this.photoID + "/comment/" + commentID).then(response => {
				if (response == null) return
				this.post_comments_cnt--;
				// Fetch comments from the server
				this.comments_data = [];
				this.page = 1;
				this.getComments();
			})
		},

		// Show or hide the comments section
		showHideComments() {
			// If comments are already shown, hide them and reset the data
			if (this.comments_shown) {
				this.comments_shown = false;
				this.comments_data = [];
				this.page = 1;
				return;
			}
			this.getComments();
		},

		// Fetch comments from the server
		getComments() {
			this.data_ended = false

			this.$axios.get("/photos/" + this.photoID + "/comment?page=" + this.page).then(response => {
					
					this.comments_shown = true;
					// If there are no more comments, set the flag
					if (response.data == null) {
						this.data_ended = true;
						return;
					}
					// Otherwise increment the start index
					else this.page++;

					// Append the comments to the array (they will be rendered)
					this.comments_data = this.comments_data.concat(response.data);
					console.log(this.comments_data);
				})
				.catch(error => {
    				console.error(error);
  				});
		},

		// Like the photo
		like() {
			this.$axios.put("/photos/" + this.photoID + "/like/" + this.$currentSession()).then(response => {
				if (response == null) return
				this.post_liked = true;
				this.post_like_cnt++;
			})
		},

		// Unlike the photo
		unlike() {
			this.$axios.delete("/photos/" + this.photoID + "/like/" + this.$currentSession()).then(response => {
				if (response == null) return
				this.post_liked = false;
				this.post_like_cnt--;
			})
		},

		// Delete the photo
		deletePhoto() {
			this.$axios.delete("/users/" + this.$currentSession() + '/photo/' + this.photoID).then(response => {
				if (response == null) return;
				this.$emit("photoDeleted");
			});
		},
	},

	created() {
		// Fetch the image from the server and display it
		console.log(this.photoID);
		this.$axios.get("/photos/" + this.photoID, { responseType: 'blob' }).then(response => {
			// Create an image element and append it to the container
			const img = document.createElement('img');

			// Set image source and css class
			img.src = URL.createObjectURL(response.data);
			img.classList.add("card-img-top");

			// Append the image to the container and disable the spinner
			this.$refs.imageContainer.appendChild(img);
			this.imageReady = true;
		});

		// Fetch the photo details from the server
		this.$axios.get("/photos/" + this.photoID + "/details").then(response => {
			if (response == null) return
			this.post_like_cnt = response.data.likesAmount;
			this.post_comments_cnt = response.data.commentsAmount;
			this.timestamp = response.data.timestamp;
			this.post_liked = response.data.isLiked;
			this.username = response.data.username;
			this.author = response.data.author;
			this.isMyPhoto = this.$currentSession() === this.author;
		});
	},
}
</script>

<template>
	<div class="card mb-5">

		<!-- Image container div -->
		<div ref="imageContainer">
			<div v-if="!imageReady" class="mt-3 mb-3">
				<LoadingSpinner :loading="!imageReady" />
			</div>
		</div>

		<div class="container">
			<div class="row">

				<!-- Userusername and timestamp -->
				<div class="col-10">
					<div class="card-body">
						<h5 @click="visitUser" class="card-title d-inline-block" style="cursor: pointer">{{ username }}</h5>
						<p class="card-text">{{ timestamp }}</p>
					</div>
				</div>

				<!-- Comment and like buttons -->
				<div class="col-2">
					<div class="card-body d-flex justify-content-end" style="display: inline-flex">
						<a v-if="isMyPhoto" @click="deletePhoto" class="btn btn-link btn-sm text-danger">
							<h5><i class="bi bi-trash"></i></h5>
						</a>
						<a @click="showHideComments">
							<h5><i class="card-title bi bi-chat-right pe-1"></i></h5>
						</a>
						<h6 class="card-text d-flex align-items-end text-muted">{{ post_comments_cnt }}</h6>
						<a v-if="!post_liked" @click="like">
							<h5><i class="card-title bi bi-suit-heart ps-2 pe-1 like-icon"></i></h5>
						</a>
						<a v-if="post_liked" @click="unlike">
							<h5><i class="card-title bi bi-heart-fill ps-2 pe-1 like-icon like-red"></i></h5>
						</a>
						<h6 class="card-text d-flex align-items-end text-muted">{{ post_like_cnt }}</h6>
						<h5></h5>
					</div>
				</div>
				</div>

			<!-- Comments section -->
			<div v-if="comments_shown">
				<div v-for="item of comments_data" class="row" v-bind:key="item.commentID">
					<div class="col-7 card-body border-top">
						<b>{{ item.username }}:</b> {{ item.commentText }}
					</div>
					<div class="col-5 card-body border-top text-end text-secondary">
						{{ item.timestamp }}
						<span v-if="isMyPhoto || item.author === this.$currentSession()" class="col-10 card-body border-top text-end">
							<button @click="deleteComment(item.commentID)" class="btn btn-link btn-sm text-danger">Delete</button>
						</span>
					</div>
				</div>

				<!-- Show more comments label -->
				<div v-if="!data_ended" class="col-12 card-body text-end pt-0 pb-1 px-0">
					<a @click="getComments" class="text-primary">Show more comments...</a>
				</div>

				<!-- New comment form -->
				<div class="row">

					<!-- Comment input -->
					<div class="col-10 card-body border-top text-end">
						<input v-model="commentMsg" type="text" class="form-control" placeholder="Commenta...">
					</div>

					<!-- Comment publish button -->
					<div class="col-1 card-body border-top text-end ps-0 d-flex">
						<button style="width: 100%" type="button" class="btn btn-primary" @click="postComment">Go</button>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<style>
.like-icon:hover {
	color: #ff0000;
}

.like-red {
	color: #ff0000;
}
</style>