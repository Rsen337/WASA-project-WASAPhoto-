<script>
    export default {
        props: ["id", "title"],
        methods: {
		    // Visit the user's profile
            showModal() {
                this.$refs['open' + this.id].click();
            },
        },
    }
</script>

<template>
    <!-- This components renders a Bootstrap modal -->
    <!-- The modal contains a title and a message -->

    <!-- Invisible button to open the modal -->
    <button :ref="'open' + id" type="button" class="btn btn-primary" style="display: none" data-bs-toggle="modal" :data-bs-target="'#' + id" />

    <!-- Modal -->
    <div class="modal fade" :id="id" tabindex="-1" :aria-labelledby="id + 'Label'" aria-hidden="true">
        <div class="modal-dialog">
            <div class="modal-content">
                <div class="modal-header">

                    <!-- Modal title -->
                    <h5 class="modal-title" :id="id + 'Label'">{{ title }}</h5>
                    <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close" @click="$emit('isClose')"></button>
                </div>

                <!-- Modal body -->
                <div class="modal-body">
                    <slot />
                </div>

                <!-- Footer with close button -->
                <div class="modal-footer">
                    <button type="button" class="btn btn-secondary" data-bs-dismiss="modal" @click="$emit('isClose')">Close</button>
                </div>
            </div>
        </div>
    </div>
</template>