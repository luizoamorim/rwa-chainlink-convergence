'use client';

import TokenizationModal from '@/components/tokenize/TokenizationModal';

export default function TokenizePage() {
	return (
		<div className="relative h-[calc(100vh-64px)] w-screen overflow-hidden flex items-center justify-center bg-gradient-to-br from-[#1E3A8A] via-[#1E40AF] to-[#0B132B]">
			{/* grid background */}

			<div
				className="absolute inset-0 opacity-[0.05]"
				style={{
					backgroundImage: 'linear-gradient(#fff 1px, transparent 1px), linear-gradient(90deg, #fff 1px, transparent 1px)',
					backgroundSize: '40px 40px',
				}}
			/>

			<TokenizationModal />
		</div>
	);
}
