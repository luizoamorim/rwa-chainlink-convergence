'use client';

import { ShieldCheck } from 'lucide-react';

type Props = {
	txHash: string;
};

export default function TokenSuccessStep({ txHash }: Props) {
	const explorerBase = process.env.NEXT_PUBLIC_EXPLORER_TX_URL;

	const explorerUrl = `${explorerBase}/${txHash}`;

	return (
		<div className="text-center">
			<ShieldCheck className="w-16 h-16 text-green-500 mx-auto mb-6" />

			<h2 className="text-xl font-semibold">Vehicle Tokenized</h2>

			<p className="text-slate-500 mt-2">Your vehicle NFT was successfully minted.</p>

			<a href={explorerUrl} target="_blank" rel="noopener noreferrer" className="block mt-6 bg-slate-100 rounded-xl p-4 text-sm break-all hover:bg-slate-200 transition">
				<div className="font-semibold mb-1">Tx Hash</div>

				{txHash}
			</a>
		</div>
	);
}
