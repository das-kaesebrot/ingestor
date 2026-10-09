import { heyApiPlugin } from '@hey-api/vite-plugin';
import tailwindcss from '@tailwindcss/vite';
import adapter from '@sveltejs/adapter-static';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	server: {
		proxy: {
			'/api/openapi.yaml': 'http://localhost:8000',
			'/api/v1': 'http://localhost:8000'
		}
	},
	plugins: [
		heyApiPlugin({
			config: {
				input: '../../openapi.yaml',
				output: 'src/lib/client',
			},
		}),
		tailwindcss(),
		sveltekit({
			compilerOptions: {
				// Force runes mode for the project, except for libraries. Can be removed in svelte 6.
				runes: ({ filename }) => filename.split(/[/\\]/).includes('node_modules') ? undefined : true
			},
			adapter: adapter({
				pages: "../static",
				assets: "../static",
				fallback: "200.html",
			})
		})
	]
});
