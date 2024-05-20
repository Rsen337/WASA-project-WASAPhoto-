<script>
// import getCurrentSession from '../services/authentication'; todo: can be removed
export default {
	data: function () {
		return {
			// The error message to display
			errormsg: null,

			loading: false,

			// Search results
			streamData: [],

			// Dynamic loading
			dataEnded: false,
			page: 1,

			// Search input
			fieldUsername: "",
		}
	},
	methods: {
		// Reset the results and fetch the new requested ones
		async query() {

			// Reset the parameters and the data
			this.dataEnded = false;
			this.streamData = [];
			this.page = 1;

			// Fetch the first batch of results
			this.loadContent();
		},

		// Fetch the search results from the server
		async loadContent() {
			this.loading = true;
			this.errormsg = null;

			// Check if the username is empty
			// and show an error message
			if (this.fieldUsername == "") {
				this.errormsg = "Please enter a username";
				this.loading = false;
				return;
			}

			console.log(this.page);
			// Fetch the results from the server
			let response = await this.$axios.get("/users?username=" + this.fieldUsername + "&page=" + this.page);

			// Errors are handled by the interceptor, which shows a modal dialog to the user and returns a null response.
			if (response == null) {
				this.loading = false
				return
			}

			// If there are no more results, set the dataEnded flag
			if (response.data.length == 0) this.dataEnded = true;

			// Otherwise, append the new results to the array
			else this.streamData = this.streamData.concat(response.data);

			// Hide the loading spinner
			this.loading = false;
		},

		// Load a new batch of results when the user scrolls to the bottom of the page
		loadMore() {
			if (this.loading || this.dataEnded) return

			this.page++;
			this.loadContent()
		},
	},
}
</script>

<template>
	<div class="mt-4">
		<div class="container">
			<div class="row justify-content-md-center">
				<div class="col-xl-6 col-lg-9">

					<h3 class="card-title border-bottom mb-4 pb-2 text-center">WASASearch</h3>

					<!-- Error message -->
					<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

					<!-- Search form -->
					<div class="form-floating mb-4">
						<input v-model="fieldUsername" @input="query" id="formUsername" class="form-control"
							placeholder="username" />
						<label class="form-label" for="formUsername">Search by username</label>
					</div>

					<!-- Search results -->
					<div id="main-content" v-for="item of streamData" v-bind:key="item.userId">
						<!-- User card (search result entry) -->
						<UserCard :user_id="item.userId" :name="item.username" :followed="false"
							:banned="false" />
					</div>

					<!-- Loading spinner -->
					<LoadingSpinner :loading="loading" /><br />

					<!-- The IntersectionObserver for dynamic loading -->
					<IntersectionObserver sentinal-name="load-more-search" @on-intersection-element="loadMore" />
				</div>
			</div>
		</div>

	</div>
</template>

<style>

</style>
