<script lang="ts">
	import type { PageProps } from './$types';
    import { page } from '$app/state';
    let token = page.params.token;

	let { data }: PageProps = $props();
</script>

<code>{data.project}</code>