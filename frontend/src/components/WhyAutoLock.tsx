import { Car, Banknote, ShieldCheck, Globe } from 'lucide-react';

export default function WhyAutoLock() {
	return (
		<section className="bg-white border-t border-slate-200">
			<div className="max-w-6xl mx-auto px-6 py-24">
				<h2 className="text-3xl font-bold text-center text-[#0B132B]">Why AutoLock DeFi</h2>

				<p className="text-center text-slate-600 mt-4 max-w-2xl mx-auto">
					Unlock the financial potential of real-world automotive assets using decentralized infrastructure.
				</p>

				<div className="grid md:grid-cols-4 gap-10 mt-16">
					{/* Market Size */}

					<div className="text-center">
						<Car size={32} className="mx-auto text-[#1E3A8A]" />

						<h3 className="text-2xl font-semibold mt-4">$3T+</h3>

						<p className="text-slate-600 mt-2">Global automotive asset market ready for tokenization</p>
					</div>

					{/* Liquidity */}

					<div className="text-center">
						<Banknote size={32} className="mx-auto text-[#1E3A8A]" />

						<h3 className="text-2xl font-semibold mt-4">Instant Liquidity</h3>

						<p className="text-slate-600 mt-2">Unlock value from vehicle ownership without selling the asset</p>
					</div>

					{/* Security */}

					<div className="text-center">
						<ShieldCheck size={32} className="mx-auto text-[#1E3A8A]" />

						<h3 className="text-2xl font-semibold mt-4">Oracle Security</h3>

						<p className="text-slate-600 mt-2">Chainlink CRE ensures decentralized verification of vehicle data</p>
					</div>

					{/* Global DeFi */}

					<div className="text-center">
						<Globe size={32} className="mx-auto text-[#1E3A8A]" />

						<h3 className="text-2xl font-semibold mt-4">Global Access</h3>

						<p className="text-slate-600 mt-2">Bring real world assets into decentralized finance markets</p>
					</div>
				</div>
			</div>
		</section>
	);
}
