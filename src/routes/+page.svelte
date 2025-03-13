<script lang="ts">
	import { goto } from '$app/navigation';
	import { getCurrentWorkout, setCurrentWorkout } from '$lib/persistence/browser/browser.svelte';

	async function startWorkout() {
		const resp = await fetch('/api/workouts', {
			method: 'POST',
		});
		const id = await resp.json();
		setCurrentWorkout(id);
		goto(`/workout/${id}`);
	}

	const currentWorkout = $state(getCurrentWorkout());
</script>

<main class="flex flex-col text-center max-w-lg space-y-2 mx-auto">
	<h1 class="text-4xl mb-8 mt-4">Bodyweight Fitness Recommended Routine</h1>

	<div>
		<button class="bg-blue-400 p-2 rounded-sm" onclick={startWorkout}>Start New Workout</button>
	</div>

	{#if currentWorkout}
		<a href={`/workout/${currentWorkout}`}>
			<button class="bg-blue-400 p-2 rounded-sm">Go to Current Workout</button>
		</a>
	{/if}
</main>
