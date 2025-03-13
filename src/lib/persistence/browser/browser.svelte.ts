import { browser } from "$app/environment";

const currentWorkoutKey = 'current-workout';

export function setCurrentWorkout(id: string) {
	if (browser) {
		window.localStorage.setItem(currentWorkoutKey, id);
	}
}

export function getCurrentWorkout(): string | null {
	if (browser) {
		return window.localStorage.getItem(currentWorkoutKey)
	}
	return null;
}
