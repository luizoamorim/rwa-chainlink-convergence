'use client';

import { motion } from 'framer-motion';
import { ShieldCheck, Network, Coins } from 'lucide-react';

const container = {
	hidden: {},
	show: {
		transition: {
			staggerChildren: 0.2,
		},
	},
};

const item = {
	hidden: { opacity: 0, y: 40 },
	show: { opacity: 1, y: 0 },
};

export default function HowItWorks() {
	return (
		<section className="bg-slate-50 border-t border-slate-200">
			<div className="max-w-6xl mx-auto px-6 py-24">
				<h2 className="text-3xl font-bold text-center text-[#0B132B]">How It Works</h2>

				<p className="text-center text-slate-600 mt-4 max-w-2xl mx-auto">
					AutoLock DeFi uses decentralized identity, oracle consensus, and smart contracts to tokenize real-world vehicle ownership.
				</p>

				<motion.div className="grid md:grid-cols-3 gap-10 mt-16" variants={container} initial="hidden" whileInView="show" viewport={{ once: true }}>
					<motion.div variants={item} className="bg-white p-8 rounded-2xl border border-slate-200">
						<ShieldCheck className="text-[#1E3A8A]" size={32} />
						<h3 className="text-xl font-semibold mt-4">Verify Identity</h3>
						<p className="text-slate-600 mt-2">Users verify their humanity using World ID.</p>
					</motion.div>

					<motion.div variants={item} className="bg-white p-8 rounded-2xl border border-slate-200">
						<Network className="text-[#1E3A8A]" size={32} />
						<h3 className="text-xl font-semibold mt-4">Oracle Consensus</h3>
						<p className="text-slate-600 mt-2">Chainlink CRE verifies vehicle registry data.</p>
					</motion.div>

					<motion.div variants={item} className="bg-white p-8 rounded-2xl border border-slate-200">
						<Coins className="text-[#1E3A8A]" size={32} />
						<h3 className="text-xl font-semibold mt-4">Mint RWA NFT</h3>
						<p className="text-slate-600 mt-2">Vehicle ownership becomes a tokenized asset.</p>
					</motion.div>
				</motion.div>
			</div>
		</section>
	);
}
