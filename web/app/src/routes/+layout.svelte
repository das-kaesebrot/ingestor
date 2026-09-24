<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
    import Head from '$lib/components/Head.svelte';

	let { children } = $props();
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
	<link href="/error" />

	<Head />
</svelte:head>
{@render children()}
