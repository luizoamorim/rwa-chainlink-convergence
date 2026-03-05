'use client';

import { useRouter } from 'next/navigation';
import { useActiveAccount } from 'thirdweb/react';
import { motion } from 'framer-motion';

export default function Hero() {
	const router = useRouter();
	const account = useActiveAccount();

	const connected = !!account?.address;

	return (
		<section className="relative overflow-hidden bg-white">
			{/* background gradient */}

			<div className="absolute inset-0 bg-gradient-to-b from-blue-50 to-transparent opacity-60" />

			{/* grid background */}

			<div
				className="absolute inset-0 opacity-[0.04]"
				style={{
					backgroundImage: 'linear-gradient(#000 1px, transparent 1px), linear-gradient(90deg, #000 1px, transparent 1px)',
					backgroundSize: '40px 40px',
				}}
			/>

			<div className="relative max-w-6xl mx-auto px-6 py-32 text-center">
				<motion.h1
					className="text-5xl font-bold text-[#0B132B] leading-tight"
					initial={{ opacity: 0, y: 40 }}
					animate={{ opacity: 1, y: 0 }}
					transition={{ duration: 0.7 }}
				>
					Tokenize Vehicles
					<span className="block text-[#1E3A8A]">Unlock DeFi Liquidity</span>
				</motion.h1>

				<motion.p className="mt-6 text-lg text-slate-600 max-w-2xl mx-auto" initial={{ opacity: 0, y: 40 }} animate={{ opacity: 1, y: 0 }} transition={{ delay: 0.2 }}>
					AutoLock DeFi connects real-world automotive assets with decentralized finance. Transform your vehicle ownership into an on-chain collateralized asset.
				</motion.p>

				<motion.div className="mt-10" initial={{ opacity: 0, y: 40 }} animate={{ opacity: 1, y: 0 }} transition={{ delay: 0.4 }}>
					<button
						disabled={!connected}
						onClick={() => router.push('/tokenize')}
						className={`px-10 py-4 rounded-xl text-lg font-semibold transition ${
							connected ? 'bg-[#1E3A8A] text-white hover:bg-[#162f6f]' : 'bg-slate-300 text-slate-500 cursor-not-allowed'
						}`}
					>
						{connected ? 'Start Tokenization' : 'Connect Wallet to Start'}
					</button>
				</motion.div>
			</div>
		</section>
	);
}
