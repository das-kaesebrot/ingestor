import type { PageLoad } from './$types';
import { BACKEND_API_ROOT } from '$lib/backend/constants'
import { Convert } from '$lib/backend/project'
import { error } from '@sveltejs/kit';

export const load: PageLoad = async ({ params, fetch }) => {
    const response = await fetch(`${BACKEND_API_ROOT}/projects/${params.token}`);

    if (!response.ok) {
        error(response.status, response.statusText);
    }

    const rawData = await response.json();
    const project = Convert.toProject(rawData)
    return { project }
};
