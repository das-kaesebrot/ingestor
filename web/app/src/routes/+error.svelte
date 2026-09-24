<script lang="ts">
	import { page } from '$app/state';

	let { status } = $props();
</script>

<svelte:head>
	<title>Error {status ?? page.status}</title>
	<meta name="robots" content="noindex,noimageindex" />
</svelte:head>

<h1>{status ?? page.status} {page.error?.message}</h1>
<span style="font-size: 10em">
	{status ?? page.status}
</span>
