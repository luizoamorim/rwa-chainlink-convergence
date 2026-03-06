'use client';

import { useEffect, useState } from 'react';
import { useActiveAccount } from 'thirdweb/react';

import TokenFormStep from './TokenFormStep';
import TokenSuccessStep from './TokenSuccessStep';
import TokenProgressStep from './TokenProgressStep';
import Pipeline from './Pipeline';

type Step = 'form' | 'progress' | 'success';

type Status = 'pending' | 'running' | 'done';

type PipelineStep = {
	label: string;
	status: Status;
};

export default function TokenizationModal() {
	const account = useActiveAccount();
	const wallet = account?.address ?? '';

	const [step, setStep] = useState<Step>('form');

	const [plate, setPlate] = useState('');
	const [renavam, setRenavam] = useState('');

	const [txHash, setTxHash] = useState('');

	const initialPipeline: PipelineStep[] = [
		{ label: 'User', status: 'pending' },
		{ label: 'World ID', status: 'pending' },
		{ label: 'Oracle', status: 'pending' },
		{ label: 'Mint NFT', status: 'pending' },
	];

	const [pipelineSteps, setPipelineSteps] = useState<PipelineStep[]>(initialPipeline);

	//////////////////////////////////////////////////////////////
	// START TOKENIZATION
	//////////////////////////////////////////////////////////////

	function startTokenization() {
		setPipelineSteps([
			{ label: 'User', status: 'running' },
			{ label: 'World ID', status: 'pending' },
			{ label: 'Oracle', status: 'pending' },
			{ label: 'Mint NFT', status: 'pending' },
		]);

		setStep('progress');
	}

	//////////////////////////////////////////////////////////////
	// STAGE HANDLER
	//////////////////////////////////////////////////////////////

	function handleStage(stage: string, tx?: string) {
		setPipelineSteps((prev) => {
			const steps = [...prev];

			switch (stage) {
				case 'verifying_identity':
					steps[0].status = 'done';
					steps[1].status = 'running';

					break;

				case 'worldid_verified':
					steps[1].status = 'done';
					steps[2].status = 'running';

					break;

				case 'checking_vehicle_registry':
					steps[1].status = 'done';
					steps[2].status = 'running';

					break;

				case 'minting_nft':
					steps[2].status = 'done';
					steps[3].status = 'running';

					break;

				case 'success':
					steps[3].status = 'done';

					if (tx) setTxHash(tx);

					setTimeout(() => {
						setStep('success');
					}, 800);

					break;
			}

			return steps;
		});
	}

	//////////////////////////////////////////////////////////////
	// WEBSOCKET
	//////////////////////////////////////////////////////////////

	useEffect(() => {
		const socket = new WebSocket('ws://localhost:8081/ws');

		socket.onmessage = (event) => {
			const data = JSON.parse(event.data);

			handleStage(data.stage, data.txHash);
		};

		socket.onerror = (err) => {
			console.warn('WS error', err);
		};

		return () => socket.close();
	}, []);

	//////////////////////////////////////////////////////////////
	// UI
	//////////////////////////////////////////////////////////////

	return (
		<div className="relative z-10 w-full max-w-xl mx-auto">
			<div className="bg-white rounded-2xl shadow-2xl p-10">
				{step === 'form' && <TokenFormStep plate={plate} setPlate={setPlate} renavam={renavam} setRenavam={setRenavam} wallet={wallet} onStart={startTokenization} />}

				{step === 'progress' && <TokenProgressStep />}

				{step === 'success' && <TokenSuccessStep txHash={txHash} />}
			</div>

			<div className="mt-10">
				<Pipeline steps={pipelineSteps} />
			</div>
		</div>
	);
}
