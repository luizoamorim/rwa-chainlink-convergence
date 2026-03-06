'use client';

import { useEffect, useState } from 'react';
import { AnimatePresence, motion } from 'framer-motion';
import { useActiveAccount } from 'thirdweb/react';

import TokenFormStep from './tokenize/TokenFormStep';
import TokenProgressStep from './tokenize/TokenProgressStep';
import TokenSuccessStep from './tokenize/TokenSuccessStep';

type Flow = 'form' | 'progress' | 'success';

export default function TokenizationForm() {
	const account = useActiveAccount();
	const address = account?.address;

	const [plate, setPlate] = useState('');
	const [renavam, setRenavam] = useState('');

	const [flow, setFlow] = useState<Flow>('form');
	const [stage, setStage] = useState<string | null>(null);
	const [txHash, setTxHash] = useState<string | null>(null);

	useEffect(() => {
		const ws = new WebSocket('ws://localhost:8081/ws');

		ws.onmessage = (event) => {
			const data = JSON.parse(event.data);

			if (data.stage) {
				setFlow('progress');
				setStage(data.stage);
			}

			if (data.txHash) {
				setTxHash(data.txHash);
				setFlow('success');
			}
		};

		return () => ws.close();
	}, []);

	return (
		<div className="min-h-[80vh] flex items-center justify-center">
			<AnimatePresence mode="wait">
				{flow === 'form' && (
					<motion.div key="form" initial={{ opacity: 0, y: 30 }} animate={{ opacity: 1, y: 0 }} exit={{ opacity: 0 }}>
						<TokenFormStep plate={plate} setPlate={setPlate} renavam={renavam} setRenavam={setRenavam} wallet={address!} />
					</motion.div>
				)}

				{flow === 'progress' && (
					<motion.div key="progress" initial={{ opacity: 0 }} animate={{ opacity: 1 }} exit={{ opacity: 0 }}>
						<TokenProgressStep stage={stage} />
					</motion.div>
				)}

				{flow === 'success' && (
					<motion.div key="success" initial={{ opacity: 0, scale: 0.95 }} animate={{ opacity: 1, scale: 1 }}>
						<TokenSuccessStep txHash={txHash} />
					</motion.div>
				)}
			</AnimatePresence>
		</div>
	);
}
