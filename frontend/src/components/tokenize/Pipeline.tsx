'use client';

import { CheckCircle2, Loader2 } from 'lucide-react';
import { motion } from 'framer-motion';

export type StepStatus = 'pending' | 'running' | 'done';

type Step = {
	label: string;
	status: StepStatus;
};

type Props = {
	steps: Step[];
};

export default function Pipeline({ steps }: Props) {
	return (
		<div className="flex items-center justify-center gap-6 text-white mt-10">
			{steps.map((step, i) => (
				<div key={step.label} className="flex items-center">
					<Node step={step} />

					{i < steps.length - 1 && <Line status={step.status} />}
				</div>
			))}
		</div>
	);
}

function Node({ step }: { step: Step }) {
	if (step.status === 'done') {
		return (
			<div className="flex flex-col items-center gap-2">
				<CheckCircle2 className="text-green-400 w-7 h-7" />
				<span className="text-xs opacity-80">{step.label}</span>
			</div>
		);
	}

	if (step.status === 'running') {
		return (
			<div className="flex flex-col items-center gap-2">
				<Loader2 className="animate-spin text-blue-400 w-7 h-7" />
				<span className="text-xs opacity-80">{step.label}</span>
			</div>
		);
	}

	return (
		<div className="flex flex-col items-center gap-2">
			<div className="w-7 h-7 rounded-full border border-white/40" />
			<span className="text-xs opacity-80">{step.label}</span>
		</div>
	);
}

function Line({ status }: { status: StepStatus }) {
	if (status === 'done') {
		return <div className="w-20 h-[2px] bg-white mx-2" />;
	}

	if (status === 'running') {
		return (
			<div className="relative w-20 h-[2px] bg-white/20 mx-2 overflow-hidden">
				<motion.div
					className="absolute h-full w-10 bg-white"
					animate={{ x: ['-30%', '120%'] }}
					transition={{
						repeat: Infinity,
						duration: 1,
						ease: 'linear',
					}}
				/>
			</div>
		);
	}

	return <div className="w-20 h-[2px] bg-white/20 mx-2" />;
}
