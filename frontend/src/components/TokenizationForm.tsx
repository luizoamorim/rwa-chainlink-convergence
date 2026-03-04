'use client';

import { useState, useEffect } from 'react';
import { useActiveAccount } from 'thirdweb/react';
import { Car } from 'lucide-react';
import VerifyButton from './VerifyButton';

export default function TokenizationForm() {
	const account = useActiveAccount();
	const address = account?.address;

	const [plate, setPlate] = useState('');
	const [renavam, setRenavam] = useState('');

	const [stage, setStage] = useState<string | null>(null);
	const [txHash, setTxHash] = useState<string | null>(null);
	const [error, setError] = useState<string | null>(null);

	useEffect(() => {
		const socket = new WebSocket('ws://localhost:8081/ws');

		socket.onmessage = (event) => {
			const data = JSON.parse(event.data);

			if (data.stage === 'received') {
				setStage('received');
			}

			if (data.stage === 'executing_cre') {
				setStage('executing');
			}

			if (data.stage === 'success') {
				setStage('success');
				setTxHash(data.txHash);
			}

			if (data.stage === 'error') {
				setStage('error');
				setError(data.error);
			}
		};

		return () => socket.close();
	}, []);

	return (
		<div className="flex justify-center mt-24 px-4">
			<div className="w-full max-w-lg bg-white rounded-3xl shadow-xl border border-slate-200 p-12">
				{/* Header */}
				<div className="flex items-center gap-4 mb-10">
					<div className="bg-blue-600 p-3 rounded-2xl shadow-md">
						<Car className="text-white" size={28} />
					</div>
					<h2 className="text-3xl font-semibold text-slate-900">Vehicle Tokenization</h2>
				</div>

				{/* Status Messages */}
				{stage === 'executing' && <div className="mb-6 p-4 bg-blue-50 text-blue-700 rounded-xl text-center">⛓ Executing on-chain transaction...</div>}

				{stage === 'success' && txHash && (
					<div className="mb-6 p-4 bg-green-50 text-green-700 rounded-xl text-center">
						🚗 Vehicle Tokenized Successfully!
						<a href={`https://sepolia.etherscan.io/tx/${txHash}`} target="_blank" className="block mt-3 underline font-medium">
							View Transaction
						</a>
					</div>
				)}

				{stage === 'error' && <div className="mb-6 p-4 bg-red-50 text-red-700 rounded-xl text-center">❌ Error: {error}</div>}

				{/* Inputs */}
				<div className="space-y-6">
					<input
						placeholder="License Plate"
						value={plate}
						onChange={(e) => setPlate(e.target.value)}
						className="w-full px-5 py-4 rounded-2xl border border-slate-300 bg-slate-50 focus:ring-2 focus:ring-blue-600"
					/>

					<input
						placeholder="RENAVAM"
						value={renavam}
						onChange={(e) => setRenavam(e.target.value)}
						className="w-full px-5 py-4 rounded-2xl border border-slate-300 bg-slate-50 focus:ring-2 focus:ring-blue-600"
					/>
				</div>

				<div className="mt-10">
					<VerifyButton plate={plate} renavam={renavam} wallet={address!} />
				</div>

				{address && <p className="text-green-600 text-center text-sm mt-6">● Wallet Connected</p>}
			</div>
		</div>
	);
}
