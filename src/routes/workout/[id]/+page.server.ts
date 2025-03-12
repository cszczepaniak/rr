import type { Actions } from './$types';

export const actions = {
	default: async (event) => {
		console.log(event.params)
		console.log(await event.request.formData())
		console.log("hey from the server")
		// TODO log the user in
	}
} satisfies Actions;
