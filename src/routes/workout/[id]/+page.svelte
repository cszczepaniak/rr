<script lang="ts">
	import IconPlay from '~icons/mdi/play-circle';
	import IconPause from '~icons/mdi/pause-circle';
	import IconNext from '~icons/mdi/skip-next-circle';
	import { page } from '$app/state';
	import { newTimer } from '$lib/timer/timer.svelte';
	import { clearCurrentWorkout } from '$lib/persistence/browser/browser.svelte';
	import { steps } from '$lib/workout/workout.ts';

	let { data } = $props();

	let stepIndex = $state(Math.min(data.workout.index, steps.length - 1));
	let currentStep = $derived(steps[stepIndex]);
	let upcomingStep = $derived(stepIndex < steps.length - 1 ? steps[stepIndex + 1] : null);
	let resting = $state(false);

	let countingDown = $state(false);
	let countdownTimer = newTimer(3);

	let mainTimer = newTimer(steps[0].duration ?? 0);

	let workoutFinished = $state(data.workout.index >= steps.length);
	$effect(() => {
		if (workoutFinished) {
			clearCurrentWorkout();
		}
	});

	async function advance() {
		mainTimer.pause();
		stepIndex++;
		if (currentStep.duration) {
			countingDown = true;
			countdownTimer.resetTo(3);
		}

		await fetch(`/api/workouts/${page.params.id}`, {
			method: 'PUT',
			body: JSON.stringify({
				index: stepIndex,
			}),
		});
	}

	async function nextStep() {
		if (upcomingStep === null) {
			workoutFinished = true;
			return;
		}

		if (currentStep.restAfter) {
			resting = true;
			mainTimer.resetTo(currentStep.restAfter);
			mainTimer.start(async () => {
				resting = false;
				await advance();
			});
			return;
		}

		await advance();
	}

	async function skipRest() {
		if (!resting) {
			return;
		}

		mainTimer.pause();
		resting = false;
		await advance();
	}
</script>

{#snippet timer(t: ReturnType<typeof newTimer>)}
	<p class="text-right text-9xl">
		<span>{t.seconds}</span>.<span>{t.fraction / 100}</span>s
	</p>
{/snippet}

<div class="flex flex-col text-center items-center">
	{#if workoutFinished}
		<p class="text-7xl mt-8 mx-2 text-center">Well done! Your workout is complete!</p>
	{:else if resting}
		<h1 class="text-7xl mt-8 mx-2 text-center">Well done! Take a rest.</h1>
		{#if upcomingStep}
			<h3 class="mt-8">Next movement:</h3>
			<h3 class="text-3xl">{upcomingStep.name}</h3>
		{/if}

		{@render timer(mainTimer)}
		<button class="text-red-400 disabled:text-red-200" onclick={skipRest}>
			<IconNext class="text-8xl" />
		</button>
	{:else}
		<h1 class="text-7xl mt-8 mx-2 text-center">{currentStep.category}</h1>
		<h3 class="mt-8">Current movement:</h3>
		<h3 class="text-3xl">{currentStep.name}</h3>

		<div class="my-4">
			{#if currentStep.reps}
				<p class="text-2xl">Do <span class="font-bold">{currentStep.reps}</span> reps!</p>
			{:else if currentStep.duration && countingDown}
				<p class="text-2xl">
					Get ready!
					{@render timer(countdownTimer)}

					<button
						onclick={() =>
							countdownTimer.start(() => {
								countingDown = false;
								mainTimer.resetTo(currentStep.duration ?? 0);
								mainTimer.start();
							})}
					>
						<IconPlay class="text-red-400 text-8xl" />
					</button>
				</p>
			{:else if currentStep.duration && !countingDown}
				<p class="text-2xl">
					Hold for <span class="font-bold">{currentStep.duration}</span> seconds!
				</p>
				{@render timer(mainTimer)}

				<div class="flex flex-col items-center mt-2">
					{#if !mainTimer.running && !mainTimer.elapsed}
						<button onclick={() => mainTimer.start()}>
							<IconPlay class="text-red-400 text-8xl" />
						</button>
					{:else}
						<button onclick={mainTimer.pause}>
							<IconPause class="text-red-400 text-8xl" />
						</button>
					{/if}
				</div>
			{/if}
		</div>

		<button
			disabled={mainTimer.running}
			class="text-red-400 disabled:text-red-200"
			onclick={nextStep}
		>
			<IconNext class="text-8xl" />
		</button>
	{/if}
</div>
