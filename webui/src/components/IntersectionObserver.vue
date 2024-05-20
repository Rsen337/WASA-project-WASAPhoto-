<script>
// This component emits an event when the sentinal element is intersecting the viewport
// (used to load more content when the user scrolls down)

// This component uses JavaScript's IntersectionObserver API

export default {
    name: 'IntersectionObserver',
    props: {
        sentinalName: {
            type: String,
            required: true,
        },
    },
    data() {
        return {
            // Whether the sentinal element is intersecting the viewport
            isIntersectingElement: false,
        }
    },
    watch: {
        // Emit an event when the sentinal element is intersecting the viewport
        isIntersectingElement: function (value) {
            if (!value) return
            this.$emit('on-intersection-element')
        },
    },
    mounted() {
        const sentinal = this.$refs[this.sentinalName]

        // Create an observer to check if the sentinal element is intersecting the viewport
        const handler = (entries) => {
            if (entries[0].isIntersecting) {
                this.isIntersectingElement = true
            }
            else {
                this.isIntersectingElement = false
            }
        }
        const observer = new window.IntersectionObserver(handler)
        observer.observe(sentinal)
    },
}
</script>

<template>
    <!-- The sentinal element -->
    <div :ref="sentinalName" class="w-full h-px relative" />
</template>