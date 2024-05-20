<script>
import IntersectionObserver from '../components/IntersectionObserver.vue';

export default {
	data: function () {
		return {
			// The profile to show
			requestedProfile: this.$route.params.user_id,

			// Loading flags
			loading: true,
			loadingError: false,

			// Profile data from the server
			user_data: [],

			// Protos data from the server
			stream_data: [],

			// Dynamic loading parameters
			data_ended: false,
			page: 1,
		};
	},
	watch: {
		'$route.params.user_id': {
			handler: function (user_id) {
				if (user_id !== null && user_id !== undefined) this.refresh()
			},
			deep: true,
			immediate: true
		}
	},
	methods: {
		async refresh() {
			if (this.$route.params.user_id == "me") {
				// If the id is "me", show the current user's profile
				this.requestedProfile = this.$currentSession()
			}
			else {
				// Otherwise, show "id"'s profile
				this.requestedProfile = this.$route.params.user_id
			}

			// Fetch profile info from the server
			this.getMainData()


			// Reset the parameters and the data
			this.data_ended = false
			this.stream_data = []
			this.page = 1

			// Fetch the first batch of posts
			this.loadContent()
		},

		// Fetch profile info from the server
		async getMainData() {
			let response = await this.$axios.get("/users/" + this.requestedProfile);
			if (response == null) {
				// An error occurred, set the error flag
				this.loading = false
				this.loadingError = true
				return
			}
			this.user_data = response.data
		},

		// Fetch photos from the server
		async loadContent() {
			this.loading = true;
			let response = await this.$axios.get("/users/" + this.requestedProfile + "/photo" + "?page=" + this.page)
			if (response.data == null) return // An error occurred. The interceptor will show a modal

			// If the server returned less elements than requested,
			// it means that there are no more photos to load
			if (response.data.length == 0)
				this.data_ended = true

			// Append the new photos to the array
			this.stream_data = response.data
			console.log(this.stream_data)

			// Disable the loading spinner
			this.loading = false
		},

		// Load more photos when the user scrolls to the bottom of the page
		loadMore() {
			// Avoid sending a request if there are no more photos
			if (this.loading || this.data_ended) return

			// Increase the start index and load more photos
			this.page += 1
			this.loadContent()
		},
	},
}
</script>

<template>
	<div class="mt-5">

		<div class="container">
			<div class="row justify-content-md-center">
				<div class="col-xl-6 col-lg-9">

					<!-- User card for profile info -->
					<UserCard :user_id="requestedProfile" :name="user_data['username']" :followed="user_data['followed']"
						:banned="user_data['banned']" :my_id="this.$currentSession" :show_new_post="true"
						@updateInfo="getMainData" @updatePosts="refresh" />

					<!-- Photos, followers and following counters -->
					<ProfileCounters :user_data="user_data" />

					<!-- Photos -->
					<div id="main-content" v-for="item of stream_data" v-bind:key="item">
						<!-- PostCard for the photo -->
						<PostCard :author="requestedProfile" :photoID="item"/>
					</div>

					<!-- Message when the end is reached -->
					<div v-if="data_ended" class="alert alert-secondary text-center" role="alert">
						You reached the end. Hooray! 👻
					</div>

					<!-- The loading spinner -->
					<LoadingSpinner :loading="loading" />

					<div class="d-flex align-items-center flex-column">
						<!-- Refresh button -->
						<button v-if="loadingError" @click="refresh" class="btn btn-secondary w-100 py-3">Retry</button>

						<!-- Load more button -->
						<button v-if="(!data_ended && !loading)" @click="loadMore" class="btn btn-secondary py-1 mb-5"
							style="border-radius: 15px">Load more</button>

						<!-- The IntersectionObserver for dynamic loading -->
						<IntersectionObserver sentinal-name="load-more-profile" @on-intersection-element="loadMore" />
					</div>
				</div>
			</div>
		</div>
	</div>
</template>