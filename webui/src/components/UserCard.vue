<script>
export default {
    props: ["user_id", "name", "followed", "banned", "show_new_post"],
    watch: {
        name: function (new_val, old_val) {
            this.username = new_val
        },
        banned: function (new_val, old_val) {
            this.user_banned = new_val
        },
        followed: function (new_val, old_val) {
            this.user_followed = new_val
        },
        user_id: function (new_val, old_val) {
            this.myself = this.$currentSession() == new_val
        },
    },
    data: function () {
        return {
            // User data
            username: this.name,
            user_followed: this.followed,
            user_banned: this.banned,

            // Whether the user is the currently logged in user
            myself: this.$currentSession() == this.user_id,

            // Whether to show the buttons to post a new photo and update the username
            show_post_form: false,
            show_username_form: false,

            // The new username
            newUsername: "",

            // The file to upload
            upload_file: null,
        }
    },
    methods: {
        // Logout the user
        logout() {
            this.$root.logout()
        },

        // Visit the user's profile
        visit() {
            this.$router.push({ path: "/profile/" + this.user_id });
        },

        // Follow the user
        follow() {
            this.$axios.put("/users/" + this.$currentSession() + "/followings/" + this.user_id)
                .then(response => {
                    if (response == null) return // the interceptors returns null if something goes bad
                    this.user_followed = true
                    this.$emit('updateInfo')
                })
        },

        // Unfollow the user
        unfollow() {
            this.$axios.delete("/users/" + this.$currentSession() + "/followings/" + this.user_id)
                .then(response => {
                    if (response == null) return
                    this.user_followed = false
                    this.$emit('updateInfo')
                })
        },

        // Ban the user
        ban() {
            this.$axios.put("/users/" + this.$currentSession() + "/banned/" + this.user_id)
                .then(response => {
                    if (response == null) return
                    this.user_banned = true
                    this.$emit('updateInfo')
                })
        },

        // Unban the user
        unban() {
            this.$axios.delete("/users/" + this.$currentSession() + "/banned/" + this.user_id)
                .then(response => {
                    if (response == null) return
                    this.user_banned = false
                    this.$emit('updateInfo')
                })
        },

        // Prepare the file to upload
        load_file(e) {
            let files = e.target.files || e.dataTransfer.files;
            if (!files.length) return
            this.upload_file = files[0]
        },

        // Upload the file
        submit_file() {
            this.$axios.post("/users/" + this.$currentSession() + "/photo", this.upload_file)
                .then(response => {
                    if (response == null) return
                    this.show_post_form = false
                    this.$emit('updatePosts')
                })
        },

        // Update the username
        updateUsername() {
            this.$axios.put("/users/" + this.$currentSession() + "/change-username", { newUsername: this.newUsername })
                .then(response => {
                    if (response == null) return
                    this.show_username_form = false
                    this.$emit('updateInfo')
                    this.username = this.newUsername
                })
        },
    },
}
</script>
<template>

    <div class="card mb-3">
        <div class="container">
            <div class="row">
                <div class="col-5">
                    <div class="card-body h-100 d-flex align-items-center">
                        <a @click="visit">
                            <h5 class="card-title mb-0 d-inline-block" style="cursor: pointer">{{ username }}</h5>
                        </a>
                    </div>
                </div>

                <!-- Whether to show one or two rows -->
                <div class="d-flex flex-column" v-bind:class="{
                                                'col-12': (myself && show_new_post),
                                                'col-sm-7': (myself && show_new_post),
                                                'col-7': !(myself && show_new_post),
                                                'align-items-end': !(myself && show_new_post),
                                                'align-items-sm-end': (myself && show_new_post),
                                            }">

                    <!-- Buttons -->
                    <div class="card-body d-flex">
                        <div v-if="!myself" class="d-flex">
                            <button v-if="!user_banned" @click="ban" type="button"
                                class="btn btn-outline-danger me-2">Ban</button>
                            <button v-if="user_banned" @click="unban" type="button"
                                class="btn btn-danger me-2">Banned</button>
                            <button v-if="!user_followed" @click="follow" type="button"
                                class="btn btn-primary">Follow</button>
                            <button v-if="user_followed" @click="unfollow" type="button"
                                class="btn btn-outline-primary">Following</button>
                        </div>

                        <!-- Users cannot follow or ban themselves -->
                        <div v-if="(myself && !show_new_post)">
                            <button disabled type="button" class="btn btn-secondary">Yourself</button>
                        </div>

                        <!-- Logout button -->
                        <div v-if="(myself && show_new_post)" class="col">
                            <button type="button" class="btn btn-outline-danger me-2" @click="logout">Logout</button>
                        </div>

                        <div class="d-flex col justify-content-end flex-row">

                            <!-- Update username button -->
                            <div v-if="(myself && show_new_post)" class="">
                                <button v-if="!show_username_form" type="button" class="btn btn-outline-secondary me-2"
                                    @click="show_username_form = true">Username</button>
                            </div>

                            <!-- Post a new photo button -->
                            <div v-if="(myself && show_new_post)" class="">
                                <button v-if="!show_post_form" type="button" class="btn btn-primary"
                                    @click="show_post_form = true">Post</button>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <!-- File input -->
            <div class="row" v-if="show_post_form">
                <div class="col-9">
                    <div class="card-body h-100 d-flex align-items-center">
                        <input @change="load_file" class="form-control form-control-lg" id="formFileLg" type="file" />
                    </div>
                </div>

                <!-- Publish button -->
                <div class="col-3">
                    <div class="card-body d-flex justify-content-end">
                        <button type="button" class="btn btn-primary btn-lg" @click="submit_file">Publish</button>
                    </div>
                </div>
            </div>

            <!-- New username form -->
            <div class="row" v-if="show_username_form">

                <!-- Username input -->
                <div class="col-10">
                    <div class="card-body h-100 d-flex align-items-center">
                        <input v-model="newUsername" class="form-control form-control-lg" id="formUsername"
                            placeholder="Your new fantastic username! 😜" />
                    </div>
                </div>

                <!-- Username update button -->
                <div class="col-2">
                    <div class="card-body d-flex justify-content-end">
                        <button type="button" class="btn btn-primary btn-lg" @click="updateUsername">Set</button>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>