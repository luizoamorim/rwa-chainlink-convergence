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
				environment="staging" // ou production se for o caso
				allow_legacy_proofs={true}
				preset={selfieCheckLegacy({ signal: 'vehicle-tokenization' })}
				handleVerify={async (result: IDKitResult) => {
					// NÃO verificar aqui!
					// Apenas retornar sucesso para o widget continuar
					return;
				}}
				onSuccess={async (result: IDKitResult) => {
					console.log('Proof generated:', result);

					// Enviar tudo direto para seu backend
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
				}}
				onError={(errorCode) => {
					console.error('IDKit error:', errorCode);
				}}
			/>
		</>
	);
}
