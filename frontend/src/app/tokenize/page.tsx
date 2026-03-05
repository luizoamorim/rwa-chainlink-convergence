'use client';

import TokenizationForm from '@/components/TokenizationForm';

export default function TokenizePage() {
	return (
		<div className="min-h-screen bg-slate-50 flex items-center justify-center px-6">
			<div className="max-w-xl w-full">
				<h1 className="text-3xl font-bold text-center mb-6">Tokenize Your Vehicle</h1>

				<p className="text-center text-gray-600 mb-10">Enter your vehicle details and verify your identity to mint the Vehicle NFT.</p>

				<TokenizationForm />
			</div>
		</div>
	);
}
