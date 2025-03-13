<script lang="ts">
	import { goto } from '$app/navigation';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	async function startWorkout() {
		const resp = await fetch('/api/workouts', {
			method: 'POST',
		});
		const id = await resp.json();
		goto(`/workout/${id}`);
	}
</script>

<h1>Bodyweight Fitness Recommended Routine</h1>
<button class="bg-blue-400" onclick={startWorkout}>Start Workout</button>

{#each data.workouts as workout}
	<pre>{JSON.stringify(workout)}</pre>
{/each}
