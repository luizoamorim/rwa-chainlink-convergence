'use client';

import { motion } from 'framer-motion';
import { User, ShieldCheck, Network, Coins } from 'lucide-react';

const steps = [
	{ icon: User, label: 'User' },
	{ icon: ShieldCheck, label: 'World ID' },
	{ icon: Network, label: 'Chainlink CRE' },
	{ icon: Coins, label: 'Vehicle NFT' },
];

export default function Architecture() {
	return (
		<section className="bg-white border-t border-slate-200">
			<div className="max-w-6xl mx-auto px-6 py-24">
				<h2 className="text-3xl font-bold text-center text-[#0B132B]">Protocol Architecture</h2>

				<div className="flex flex-wrap justify-center items-center gap-10 mt-16">
					{steps.map((step, i) => {
						const Icon = step.icon;

						return (
							<motion.div
								key={i}
								initial={{ opacity: 0, y: 30 }}
								whileInView={{ opacity: 1, y: 0 }}
								transition={{ delay: i * 0.2 }}
								viewport={{ once: true }}
								className="flex flex-col items-center"
							>
								<div className="bg-slate-50 p-6 rounded-2xl border border-slate-200 text-[#1E3A8A]">
									<Icon size={32} />
								</div>

								<p className="mt-3 font-medium text-slate-700">{step.label}</p>
							</motion.div>
						);
					})}
				</div>
			</div>
		</section>
	);
}
