import type { RequestHandler } from './$types';
import { list } from "@vercel/blob";
import { BLOB_READ_WRITE_TOKEN } from '$env/static/private'

export const GET: RequestHandler = async ({ }) => {
	const workoutMetas = await list({
		prefix: 'workouts/',
		token: BLOB_READ_WRITE_TOKEN,
	});

	const workouts = await Promise.all(workoutMetas.blobs.map(async b => {
		const workoutResp = await fetch(b.downloadUrl);
		return await workoutResp.json()
	}))

	return new Response(JSON.stringify(workouts));
};
