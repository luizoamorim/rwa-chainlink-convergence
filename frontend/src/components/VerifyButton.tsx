'use client';

import { IDKitRequestWidget, selfieCheckLegacy, type RpContext, type IDKitResult } from '@worldcoin/idkit';
import { useEffect, useState } from 'react';

const APP_ID = process.env.NEXT_PUBLIC_WLD_APP_ID as `app_${string}`;
const RP_ID = process.env.NEXT_PUBLIC_WLD_RP_ID!;
const ACTION = 'tokenizevehicle';

type Props = {
	plate: string;
	renavam: string;
	wallet?: string;
};

export default function VerifyButton({ plate, renavam, wallet }: Props) {
	const [open, setOpen] = useState(false);
	const [rpContext, setRpContext] = useState<RpContext | null>(null);

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

	return (
		<>
			<button onClick={() => setOpen(true)} className="w-full mt-6 py-3 bg-blue-600 text-white rounded-xl font-bold hover:bg-blue-700 transition">
				Verify with World ID
			</button>

			<IDKitRequestWidget
				open={open}
				onOpenChange={setOpen}
				app_id={APP_ID}
				action={ACTION}
				rp_context={rpContext}
				allow_legacy_proofs={true}
				preset={selfieCheckLegacy({ signal: 'vehicle-tokenization' })}
				handleVerify={async (result: IDKitResult) => {
					// 1️⃣ Verify proof server-side
					const verifyResponse = await fetch('/api/verify-proof', {
						method: 'POST',
						headers: { 'content-type': 'application/json' },
						body: JSON.stringify({
							rp_id: RP_ID,
							idkitResponse: result,
						}),
					});

					if (!verifyResponse.ok) {
						throw new Error('World ID verification failed');
					}

					// 2️⃣ Trigger tokenization
					const tokenizeResponse = await fetch('/api/tokenize', {
						method: 'POST',
						headers: { 'content-type': 'application/json' },
						body: JSON.stringify({
							plate,
							renavam,
							wallet,
							proof: result,
						}),
					});

					if (!tokenizeResponse.ok) {
						const error = await tokenizeResponse.text();
						throw new Error(error);
					}
				}}
			/>
		</>
	);
}
