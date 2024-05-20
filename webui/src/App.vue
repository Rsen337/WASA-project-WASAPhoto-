<script>
export default {
	props: ["user_id", "name", "date", "comments", "likes", "photo_id", "liked"],
	data: function () {
		return {
			// Data for the modal
			modalTitle: "Modal Title",
			modalMsg: "Modal Message",

			// Whether the user is logged in
			logged_in: true,
		}
	},
	methods: {
		// Function to show a modal
		// can be called by any view or component
		// title: title of the modal
		// message: message to show in the modal
		showModal(title, message) {
			// Set the modal data
			this.modalTitle = title
			this.modalMsg = message

			// Show the modal
			this.$refs.errModal.showModal()
		},

		// Sets the login status to true
		// to show the navigation buttons
		setLoggedIn() {
			this.logged_in = true
		},

		// Disconnects the current logged in user
		logout() {
			localStorage.removeItem("token")
            sessionStorage.removeItem("token")
			this.logged_in = false
            this.$router.push({ path: "/login" })
		}
	},

	// Called when the root view is mounted
	mounted() {
		// Check if the user is already logged in
		this.$axiosUpdate()

		// Configure axios interceptors
		this.$axios.interceptors.response.use(response => {
			// Leave response as is
			return response;
		}, error => {
			if (error.response.status != 0) {
				// If the response is 401, redirect to /login
				if (error.response.status === 401) {
					this.$router.push({ path: '/login' })
					this.logged_in = false;
					return;
				}
				
				// Show the error message from the server in a modal
				this.showModal("Error " + error.response.status, error.response.data);
				return;
			}
			// Show the error message from axios in a modal
			this.showModal("Error", error.toString());
			return;
		});
	}
}
</script>

<template>
	<!-- Modal to show error messages -->
	<Modal ref="errModal" id="errorModal" :title="modalTitle">
		{{ modalMsg }}
	</Modal>

	<div class="container-fluid">
		<div class="row">
			<main>
				<!-- The view is rendered here -->
				<RouterView />
				<div v-if="logged_in" class="mb-5 pb-3"></div> <!-- Empty div to avoid hiding items under the navbar. todo: find a better way to do this -->
			</main>

			<!-- Bottom navigation buttons -->
			<nav v-if="logged_in" id="global-nav" class="navbar fixed-bottom navbar-light bg-light">
				<div class="collapse navbar-collapse" id="navbarNav"></div>
				<RouterLink to="/" class="col-4 text-center">
					<i class="bi bi-house text-dark" style="font-size: 2em"></i>
					<span class="nav-label router-button">Home</span>
				</RouterLink>
				<RouterLink to="/search" class="col-4 text-center">
					<i class="bi bi-search text-dark" style="font-size: 2em"></i>
					<span class="nav-label router-button">Search</span>
				</RouterLink>
				<RouterLink to="/profile/me" class="col-4 text-center">
					<i class="bi bi-person text-dark" style="font-size: 2em"></i>
					<span class="nav-label router-button">Profile</span>
				</RouterLink>
			</nav>
		</div>
	</div>
</template>

<style>
/* Make the active navigation button a little bit bigger */
#global-nav a.router-link-active {
	font-size: 1.2em
}
</style>
