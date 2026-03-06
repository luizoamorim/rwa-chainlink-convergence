'use client';

import { Car, Loader2 } from 'lucide-react';
import { motion } from 'framer-motion';

export default function TokenProgressStep() {
	return (
		<div className="flex flex-col items-center justify-center text-center">
			{/* Animated Car Icon */}

			<motion.div
				className="mb-6"
				animate={{
					y: [0, -4, 0],
				}}
				transition={{
					duration: 1.5,
					repeat: Infinity,
					ease: 'easeInOut',
				}}
			>
				<Car size={48} className="text-[#1E3A8A]" />
			</motion.div>

			{/* Spinner */}

			<Loader2 size={22} className="animate-spin text-[#1E3A8A] mb-4" />

			{/* Title */}

			<h2 className="text-xl font-semibold text-slate-800">Tokenizing Vehicle</h2>

			{/* Description */}

			<p className="text-slate-500 mt-2 max-w-sm">Fetching registry data, validating ownership and minting your vehicle NFT.</p>
		</div>
	);
}
