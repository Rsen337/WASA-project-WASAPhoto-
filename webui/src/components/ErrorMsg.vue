<script>
export default {
	// The error message to display
	props: ['msg']
}
</script>

<template>
	<!-- This component renders an error message -->
	<div class="alert alert-danger" role="alert">
		{{ msg }}
	</div>
</template>

<style>
</style>
