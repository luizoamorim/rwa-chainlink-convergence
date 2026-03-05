'use client';

import { useState } from 'react';
import { useActiveAccount } from 'thirdweb/react';
import { Car } from 'lucide-react';
import VerifyButton from './VerifyButton';

export default function TokenizationForm() {
	const account = useActiveAccount();
	const address = account?.address;

	const [plate, setPlate] = useState('');
	const [renavam, setRenavam] = useState('');

	return (
		<div className="flex justify-center">
			<div className="w-full max-w-xl bg-white rounded-3xl border border-slate-200 shadow-lg p-12">
				<div className="flex items-center gap-4 mb-10">
					<div className="bg-[#1E3A8A] p-3 rounded-2xl">
						<Car className="text-white" size={28} />
					</div>

					<h2 className="text-3xl font-semibold">Vehicle Tokenization</h2>
				</div>

				<div className="space-y-6">
					<input
						placeholder="License Plate"
						value={plate}
						onChange={(e) => setPlate(e.target.value)}
						className="w-full px-5 py-4 rounded-xl border border-slate-300 focus:ring-2 focus:ring-[#1E3A8A]"
					/>

					<input
						placeholder="RENAVAM"
						value={renavam}
						onChange={(e) => setRenavam(e.target.value)}
						className="w-full px-5 py-4 rounded-xl border border-slate-300 focus:ring-2 focus:ring-[#1E3A8A]"
					/>
				</div>

				<div className="mt-10">
					<VerifyButton plate={plate} renavam={renavam} wallet={address!} />
				</div>
			</div>
		</div>
	);
}
