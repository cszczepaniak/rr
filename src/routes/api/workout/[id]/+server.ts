import type { RequestHandler } from './$types';
import { put } from "@vercel/blob";
import { BLOB_READ_WRITE_TOKEN } from '$env/static/private'

type PersistedWorkout = {
	completedIndex: number,
}

type PostRequest = {
	index: number,
}

export const POST: RequestHandler = async (req) => {
	const reqJSON: PostRequest = await req.request.json();

	const blobName = `workouts/${req.params.id}.json`
	const persisted: PersistedWorkout = { completedIndex: reqJSON.index };

	await put(blobName, JSON.stringify(persisted), {
		access: 'public',
		token: BLOB_READ_WRITE_TOKEN,
		addRandomSuffix: false
	});

	return new Response(JSON.stringify("hey"));
};
