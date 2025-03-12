<script lang="ts">
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
</script>

<h1>Bodyweight Fitness Recommended Routine</h1>
<a href="/workout/123">
	<button class="bg-blue-400">Start Workout</button>
</a>

{#each data.workouts as workout}
	<pre>{JSON.stringify(workout)}</pre>
{/each}
