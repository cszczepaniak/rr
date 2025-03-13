import type { RequestHandler } from './$types';
import { BlobNotFoundError, head, put } from "@vercel/blob";
import { BLOB_READ_WRITE_TOKEN } from '$env/static/private'

type PersistedWorkout = {
	index: number,
}

type PostRequest = {
	index: number,
}

export const PUT: RequestHandler = async (req) => {
	const reqJSON: PostRequest = await req.request.json();


	const blobName = `workouts/${req.params.id}.json`
	const persisted: PersistedWorkout = { index: reqJSON.index };
	console.log("updating", blobName, "to", persisted);

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

	return new Response();
};

export const GET: RequestHandler = async (req) => {
	const blobName = `workouts/${req.params.id}.json`

	try {
		const details = await head(blobName, { token: BLOB_READ_WRITE_TOKEN })
		const resp = await fetch(details.url);

		const parsed = await resp.json();
		console.log(blobName, parsed);

		return new Response(JSON.stringify(parsed));
	} catch (e) {
		if (e instanceof BlobNotFoundError) {
			// If we didn't find a blob, assume the workout is at the beginning.
			return new Response(JSON.stringify({ index: 0 }));
		}

		// Otherwise we'd like to know about this error!
		console.log(e);
		throw e;
	}
};
