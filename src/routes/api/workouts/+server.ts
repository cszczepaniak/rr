import type { RequestHandler } from './$types';
import { list, put } from "@vercel/blob";
import { BLOB_READ_WRITE_TOKEN } from '$env/static/private'
import type { Workout } from '$lib/persistence/workout/workout';

export const POST: RequestHandler = async () => {
	const id = crypto.randomUUID();

	const blobName = `workouts/${id}.json`
	const persisted: Workout = { id, index: 0 };

	console.log("creating", blobName, "with id", id);

	try {
		await put(blobName, JSON.stringify(persisted), {
			access: 'public',
			token: BLOB_READ_WRITE_TOKEN,
			addRandomSuffix: false,
			cacheControlMaxAge: 0,
		});
	} catch (e) {
		console.log(e)
	}

	return new Response(JSON.stringify(id));
}


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
