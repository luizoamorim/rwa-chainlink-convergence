'use client';

import { ShieldCheck } from 'lucide-react';

type Props = {
	txHash: string;
};

export default function TokenSuccessStep({ txHash }: Props) {
	return (
		<div className="text-center space-y-6">
			<ShieldCheck className="w-16 h-16 text-green-500 mx-auto" />

			<h2 className="text-xl font-semibold">Vehicle Tokenized</h2>

			<p className="text-sm text-gray-500">Your vehicle NFT was successfully minted.</p>

			{txHash && (
				<div className="bg-gray-100 rounded-xl p-4 text-sm break-all">
					<span className="font-medium">Tx Hash:</span>
					<br />
					{txHash}
				</div>
			)}
		</div>
	);
}
