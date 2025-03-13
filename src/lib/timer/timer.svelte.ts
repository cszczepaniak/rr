export function newTimer(initialDuration: number) {
	let timerHandle = $state(0);
	let milliseconds = $state(initialDuration * 1000);
	let running = $state(false);
	let seconds = $derived(Math.floor(milliseconds / 1000));
	let fraction = $derived(milliseconds % 1000);
	let elapsed = $derived(milliseconds <= 0);

	let tickMilliseconds = 100;

	return {
		resetTo: (seconds: number) => { milliseconds = seconds * 1000 },
		start: (onDone?: () => void) => {
			if (running) {
				return;
			}

			running = true;
			timerHandle = setInterval(() => {
				console.log(`tick: ${milliseconds}`);
				milliseconds -= tickMilliseconds;
				if (milliseconds == 0) {
					running = false;
					clearInterval(timerHandle);
					if (onDone) {
						onDone();
					}
				}
			}, tickMilliseconds)
		},
		pause: () => {
			if (!running) {
				return;
			}

			running = false;
			clearInterval(timerHandle);
			timerHandle = 0;
		},
		get elapsed() {
			return elapsed;
		},
		get running() {
			return running;
		},
		get seconds() {
			return seconds;
		},
		get fraction() {
			return fraction;
		},
	}
}
