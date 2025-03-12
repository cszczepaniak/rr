import type { PageLoad } from './$types'

export const load: PageLoad = async (event) => {
	const workoutResp = await event.fetch('/api/workouts');
	return {
		workouts: await workoutResp.json(),
	};
};
