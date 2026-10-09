<script lang="ts">
  import { onMount } from 'svelte';
  import swaggerCss from 'swagger-ui-dist/swagger-ui.css?url';

  let { url = '/openapi.yaml', spec = undefined }: { url?: string; spec?: object } = $props();

  let container: HTMLDivElement;

  onMount(async () => {
    const mod: any = await import('swagger-ui-dist/swagger-ui-bundle.js');
    const SwaggerUIBundle = mod.default;

    SwaggerUIBundle({
      domNode: container,
      ...(spec ? { spec } : { url }),
      deepLinking: true
    });
  });
</script>

<svelte:head>
  <link rel="stylesheet" href={swaggerCss} />
</svelte:head>

<div bind:this={container}></div>