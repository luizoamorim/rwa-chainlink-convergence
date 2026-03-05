'use client';

import { IDKitRequestWidget, type IDKitResult, type RpContext, selfieCheckLegacy } from '@worldcoin/idkit';
import { useEffect, useState } from 'react';

const APP_ID = process.env.NEXT_PUBLIC_WLD_APP_ID as `app_${string}`;
const RP_ID = process.env.NEXT_PUBLIC_WLD_RP_ID!;
const ACTION = process.env.NEXT_PUBLIC_WLD_ACTION!;

type Props = {
	plate: string;
	renavam: string;
	wallet: string;
};

export default function VerifyButton({ plate, renavam, wallet }: Props) {
	const [open, setOpen] = useState(false);
	const [rpContext, setRpContext] = useState<RpContext | null>(null);
	const [isProcessing, setIsProcessing] = useState(false);

	useEffect(() => {
		async function fetchSignature() {
			const res = await fetch('/api/rp-signature', {
				method: 'POST',
				headers: { 'content-type': 'application/json' },
				body: JSON.stringify({ action: ACTION }),
			});

			const data = await res.json();

			setRpContext({
				rp_id: RP_ID,
				nonce: data.nonce,
				created_at: data.created_at,
				expires_at: data.expires_at,
				signature: data.sig,
			});
		}

		fetchSignature();
	}, []);

	if (!rpContext) return null;

	const disabled = isProcessing || !plate || !renavam || !wallet;

	return (
		<>
			<button
				disabled={disabled}
				onClick={() => setOpen(true)}
				className={`w-full mt-6 py-3 rounded-xl font-bold transition ${
					disabled ? 'bg-slate-400 text-white cursor-not-allowed' : 'bg-blue-600 text-white hover:bg-blue-700'
				}`}
			>
				{isProcessing ? 'Processing...' : 'Verify with World ID'}
			</button>

			<IDKitRequestWidget
				open={open}
				onOpenChange={setOpen}
				app_id={APP_ID}
				action={ACTION}
				rp_context={rpContext}
				environment="staging"
				allow_legacy_proofs={true}
				preset={selfieCheckLegacy({ signal: 'vehicle-tokenization' })}
				handleVerify={async () => {
					// Let widget continue
					return;
				}}
				onSuccess={async (result: IDKitResult) => {
					setIsProcessing(true);

					try {
						await fetch('/api/tokenize', {
							method: 'POST',
							headers: { 'content-type': 'application/json' },
							body: JSON.stringify({
								plate,
								renavam,
								wallet,
								proof: result,
							}),
						});
					} finally {
						setIsProcessing(false);
					}
				}}
				onError={(errorCode) => {
					console.error('IDKit error:', errorCode);
				}}
			/>
		</>
	);
}
