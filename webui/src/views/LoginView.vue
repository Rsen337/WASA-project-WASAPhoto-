<script>
export default {
    data: function () {
        return {
            // The error message to display
            errormsg: null,

            // Loading spinner state
            loading: false,

            // Form inputs
            field_username: "",
            rememberLogin: false,
        };
    },
    methods: {
        // Send the login request to the server
        // if the login is successful, the token is saved
        // and the user is redirected to the previous page
        async login() {
            this.loading = true;
            this.errormsg = null;

            // Send the login request
            let response = await this.$axios.post("/session", {
                token: this.field_username,
            });


            // Errors are handled by the interceptor, which shows a modal dialog to the user and returns a null response.
            if (response == null) {
				this.loading = false
				return
			}

            // If the login is successful, save the token and redirect to the previous page
            if (response.status == 201 || response.status == 200) {
                // Save the token in the local storage if the user wants to be remembered
                localStorage.setItem("username", this.field_username);
                if (this.rememberLogin) {
                    localStorage.setItem("token", response.data["identifier"])
                    sessionStorage.removeItem("token");
                }
                // Else save the token in the session storage
                else {
                    sessionStorage.setItem("token", response.data["identifier"]);
                    localStorage.removeItem("token");
                }
                // Tell the root view to enable the navbar
                this.$root.setLoggedIn();
                // Update the header
                this.$axiosUpdate();

                // Go back to the previous page
                this.$router.go(-1);
            }
            else {
                // Login failed, show the error message
                this.errormsg = response.data["error"];
            }
            // Disable the loading spinner
            this.loading = false;
        },
    },
}
</script>

<template>
    <!-- Login form centered in the page -->
    <div class="vh-100 container py-5 h-100">
        <div class="row d-flex justify-content-center align-items-center h-100">
            <!--<div class="col-sm"><h2>* immagina un logo carino *</h2></div>-->
            <div class="col-12 col-md-8 col-lg-6 col-xl-5">
                <div class="card" style="border-radius: 1rem">
                    <div class="card-body p-4">

                        <h1 class="h2 pb-4 text-center">WASAPhoto</h1>

                        <form>
                            <!-- Email input -->
                            <div class="form-floating mb-4">
                                <input v-model="field_username" type="text" id="formUsername" class="form-control"
                                    placeholder="username" />
                                <label class="form-label" for="formUsername">Username</label>
                            </div>

                            <!-- Password input -->
                            <div class="form-floating mb-4">
                                <input style="display: none" disabled type="password" id="formPassword"
                                    class="form-control" placeholder="gattina12" />
                                <label style="display: none" class="form-label" for="formPassword">Password</label>
                            </div>

                            <!-- 2 column grid layout for inline styling -->
                            <div class="row mb-4">
                                <div class="col d-flex justify-content-center">
                                    <!-- Checkbox -->
                                    <div class="form-check">
                                        <input v-model="rememberLogin" class="form-check-input" type="checkbox" value=""
                                            id="form2Example31" />
                                        <label class="form-check-label" for="form2Example31">Remember me</label>
                                    </div>
                                </div>
                            </div>

                            <!-- Submit button -->
                            <button style="width: 100%" type="button" class="btn btn-primary btn-block mb-4"
                                @click="login">Sign in</button>
                            <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
                            <LoadingSpinner :loading="loading" />
                            <i class="text-center text-secondary d-flex flex-column">share your special moments!</i>
                        </form>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<style>

</style>
