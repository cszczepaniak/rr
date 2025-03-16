type Progression = {
	category: string;
	exercises: Step[];
}

type Step = {
	name: string;
	reps?: number;
	duration?: number;
	restAfter?: number;
};

export const progressions: Progression[] = [
	{
		category: 'Warmup',
		exercises: [
			{
				name: 'Shoulder Circles',
				reps: 10,
			},
		],
	},
	{
		category: 'Warmup',
		exercises: [
			{
				name: 'Scapular Shrugs',
				reps: 10,
			},
		],
	},
	{
		category: 'Warmup',
		exercises: [
			{
				name: 'Cat/Camel',
				reps: 10,
			},
		],
	},
	{
		category: 'Warmup',
		exercises: [
			{
				name: 'Band Pulldowns',
				reps: 10,
			},
		],
	},
	{
		category: 'Warmup',
		exercises: [
			{
				name: 'Band Dislocates',
				reps: 10,
			},
		],
	},
	{
		category: 'Warmup',
		exercises: [
			{
				name: 'Wrist Mobility',
			},
		],
	},
	{
		category: 'Warmup',
		exercises: [
			{
				name: 'Hamstring Stretch (1 of 2)',
				duration: 30,
			},
		],
	},
	{
		category: 'Warmup',
		exercises: [
			{
				name: 'Hamstring Stretch (2 of 2)',
				duration: 30,
			},
		],
	},
	{
		category: 'Skill Work',
		exercises: [
			{
				name: 'Parallel Bar Support',
				duration: 30,
			},
			{
				name: 'Ring Support',
				duration: 30,
			},
			{
				name: 'Ring RTO Support',
				duration: 30,
			},
		],
	},
	{
		category: 'Skill Work',
		exercises: [
			{
				name: 'Wall Handstand',
				duration: 30,
			},
			{
				name: 'Handstand',
				duration: 30,
			},
		],
	},
]

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

export const steps: Step[] = [
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
