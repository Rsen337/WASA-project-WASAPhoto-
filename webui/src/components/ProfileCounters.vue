<script>
export default {
    props: ["user_data"],

    data: function () {
        return {
            modal_data: [],

            data_type: "followers",

            // Dynamic loading parameters
            data_ended: false,
            page: 1,
            loading: false,
        };
    },
    methods: {
        // Visit the profile of the user with the given id
        visit(user_id) {
            this.$router.push({ path: "/profile/" + user_id })
        },

        // Reset the current data and fetch the first batch of users
        // then show the modal
        async loadData(type) {
            // Reset the parameters and the users array
            this.data_type = type
            this.data_ended = false
            this.modal_data = []

            // Fetch the first batch of users
            let status = await this.loadContent()

            // Show the modal if the request was successful
            if (status) this.$refs["mymodal"].showModal()
            // If the request fails, the interceptor will show the error modal
        },

        // Fetch users from the server
        async loadContent() {
            // Fetch followers / following from the server
            // uses /followers and /following endpoints
            if (this.data_type == "followers")
                var response = await this.$axios.get("/users/" + this.user_data["userId"] + "/followers?page=" + this.page)
            else if (this.data_type == "followings")
                var response = await this.$axios.get("/users/" + this.user_data["userId"] + "/followings?page=" + this.page)

            if (response.data == null) return false // An error occurred. The interceptor will show a modal

            // If the server returned less elements than requested,
            // it means that there are no more photos to load
            if (response.data.length == 0)
                this.data_ended = true

            // Append the new photos to the array
            this.modal_data = this.modal_data.concat(response.data)
            return true
        },

        // Load more users when the user scrolls to the bottom
        loadMore() {
            // Avoid sending a request if there are no more photos
            if (this.loading || this.data_ended) return

            // Increase the start index and load more photos
            this.page++
            this.loadContent()
        },
    },
}
</script>

<template>

    <!-- Modal to show the followers / following -->
    <Modal ref="mymodal" id="userModal" :title="data_type" @isClose="this.page=1">
        <ul>
            <li v-for="item in modal_data" :key="item.userId" class="mb-2" style="cursor: pointer"
                @click="visit(item.userId); this.page=1" data-bs-dismiss="modal">
                <h5>{{ item.username }}</h5>
            </li>
            <IntersectionObserver sentinal-name="load-more-users" @on-intersection-element="loadMore" />
        </ul>
    </Modal>

    <!-- Profile counters -->
    <div class="row text-center mt-2 mb-3">

        <!-- Photos counter -->
        <div class="col-4" style="border-right: 1px">
            <h3>{{ user_data["photos"] }}</h3>
            <h6>Photos</h6>
        </div>

        <!-- Followers counter -->
        <div class="col-4" @click="loadData('followers')" style="cursor: pointer">
            <h3>{{ user_data["followers"] }}</h3>
            <h6>Followers</h6>
        </div>

        <!-- Following counter -->
        <div class="col-4" @click="loadData('followings')" style="cursor: pointer">
            <h3>{{ user_data["followings"] }}</h3>
            <h6>Following</h6>
        </div>
    </div>

</template>
