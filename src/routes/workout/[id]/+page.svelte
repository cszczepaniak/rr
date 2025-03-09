<script lang="ts">
	type Step = {
		category: string;
		name: string;
		reps?: number;
		duration?: number;
		restAfter?: number;
	};

	function repeat(s: Step[], times: number): Step[] {
		const result: Step[] = [];
		for (let i = 0; i < times; i++) {
			result.push(...s);
		}
		return result;
	}

	function category(name: string, steps: Omit<Step, 'category'>[]): Step[] {
		return steps.map((s) => ({ ...s, category: name }));
	}

	function fixedRestCategory(
		name: string,
		rest: number,
		steps: Omit<Omit<Step, 'category'>, 'restAfter'>[],
	): Step[] {
		return steps.map((s) => ({ ...s, category: name, restAfter: rest }));
	}

	const steps: Step[] = [
		...category('Warmup', [
			{
				name: 'Shoulder Circles',
				reps: 10,
			},
			{
				name: 'Scapular Shrugs',
				reps: 10,
			},
			{
				name: 'Cat/Camel',
				reps: 10,
			},
			{
				name: 'Band Pulldowns',
				reps: 10,
			},
			{
				name: 'Band Dislocates',
				reps: 10,
			},
			{
				name: 'Wrist Mobility',
			},
			{
				name: 'Hamstring Stretch (1 of 2)',
				duration: 30,
			},
			{
				name: 'Hamstring Stretch (2 of 2)',
				duration: 30,
			},
		]),
		...repeat(
			fixedRestCategory('Skill Work', 60, [
				{
					name: 'Parallel Bar Support',
					duration: 30,
				},
				{
					name: 'Wall Handstand',
					duration: 30,
				},
			]),
			3,
		),
		...repeat(
			fixedRestCategory('Strength Work (Set 1)', 90, [
				{
					name: 'Pseudo-planche Pushup',
					reps: 8,
				},
				{
					name: 'Wide Rows',
					reps: 8,
				},
			]),
			3,
		),
		...repeat(
			fixedRestCategory('Strength Work (Set 2)', 90, [
				{
					name: 'L-Sit (Foot-supported)',
					duration: 30,
				},
				{
					name: 'One-leg stepups',
					reps: 8,
				},
			]),
			3,
		),
		...repeat(
			fixedRestCategory('Strength Work (Set 3)', 90, [
				{
					name: 'Pullups',
					reps: 8,
				},
				{
					name: 'Parallel Bar Dips',
					reps: 8,
				},
			]),
			3,
		),
		...repeat(
			category('Bodyline Drills', [
				{
					name: 'Plank (Elbows)',
					duration: 30,
					restAfter: 60,
				},
				{
					name: 'Side Plank (L Elbow)',
					duration: 30,
				},
				{
					name: 'Side Plank (R Elbow)',
					duration: 30,
					restAfter: 60,
				},
				{
					name: 'Hollow Hold (Feet Up)',
					duration: 30,
					restAfter: 60,
				},
				{
					name: 'Superman',
					duration: 30,
					restAfter: 60,
				},
			]),
			3,
		),
	];

	let stepIndex = $state(0);
	let currentStep = $derived(steps[stepIndex]);
	let resting = $state(false);

	let timerHandle = $state(0);
	let timer = $state(steps[0].duration ?? 0);
	let timerRunning = $state(false);

	let workoutFinished = $state(false);

	function advance() {
		stopTimer();
		stepIndex++;
		if (currentStep.duration) {
			timer = currentStep.duration;
		}
	}

	function nextStep() {
		if (stepIndex >= steps.length - 1) {
			workoutFinished = true;
			return;
		}

		if (currentStep.restAfter) {
			resting = true;
			timer = currentStep.restAfter;
			startTimer(() => {
				resting = false;
				advance();
			});
			return;
		}

		advance();
	}

	function skipRest() {
		if (!resting) {
			return;
		}

		stopTimer();
		resting = false;
		advance();
	}

	function startTimer(onDone?: () => void) {
		if (timerRunning) {
			return;
		}

		timerRunning = true;

		timerHandle = setInterval(() => {
			timer--;
			if (timer == 0) {
				stopTimer();
				if (onDone) {
					onDone();
				}
			}
		}, 1000);
	}

	function stopTimer() {
		if (!timerRunning) {
			return;
		}

		timerRunning = false;
		clearInterval(timerHandle);
	}
</script>

<div class="flex flex-col items-center">
	{#if workoutFinished}
		<p class="text-7xl mt-8">Well done! Your workout is complete!</p>
	{:else}
		<h1 class="text-7xl mt-8">{currentStep.category}</h1>
		<h3 class="mt-8">Current movement:</h3>
		<h3 class="text-3xl">{currentStep.name}</h3>

		{#if resting}
			<p>Well done! Take a rest.</p>
			<p class="text-center text-9xl">{timer}s</p>
			<button
				class="bg-red-400 disabled:bg-red-200 disabled:text-gray-400 p-1 rounded-sm"
				onclick={skipRest}>Skip Rest</button
			>
		{:else}
			<div class="my-4">
				{#if currentStep.reps}
					<p class="text-2xl">Do <span class="font-bold">{currentStep.reps}</span> reps!</p>
				{:else if currentStep.duration}
					<p class="text-2xl">
						Hold for <span class="font-bold">{currentStep.duration}</span> seconds!
					</p>
					<p class="text-center text-9xl">{timer}s</p>

					<div class="flex flex-col items-center mt-2">
						{#if !timerRunning}
							<button class="bg-red-400 p-1 rounded-sm" onclick={() => startTimer()}>Start</button>
						{:else}
							<button class="bg-red-400 p-1 rounded-sm" onclick={stopTimer}>Stop</button>
						{/if}
					</div>
				{/if}
			</div>

			<button
				disabled={timerRunning}
				class="bg-red-400 disabled:bg-red-200 disabled:text-gray-400 p-1 rounded-sm"
				onclick={nextStep}>Done</button
			>
		{/if}
	{/if}
</div>
