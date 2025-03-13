import type { PageLoad } from './$types'

export const ssr = false;

export const load: PageLoad = async (event) => {
	const workoutResp = await event.fetch(`/api/workouts/${event.params.id}`);

	return {
		workout: await workoutResp.json(),
	};
};
