'use client';

import { Car } from 'lucide-react';
import VerifyButton from '../VerifyButton';

type Props = {
	plate: string;
	setPlate: (v: string) => void;
	renavam: string;
	setRenavam: (v: string) => void;
	wallet: string;
	onStart: () => void;
};

export default function TokenFormStep({ plate, setPlate, renavam, setRenavam, wallet, onStart }: Props) {
	return (
		<div>
			<div className="flex items-center gap-4 mb-8">
				<div className="bg-[#1E3A8A] p-3 rounded-xl">
					<Car className="text-white" size={26} />
				</div>

				<h2 className="text-2xl font-semibold text-[#0B132B]">Tokenize Vehicle</h2>
			</div>

			<div className="space-y-4">
				<input
					placeholder="License Plate"
					value={plate}
					onChange={(e) => setPlate(e.target.value)}
					className="w-full px-4 py-3 rounded-lg border border-slate-300 focus:ring-2 focus:ring-[#1E3A8A]"
				/>

				<input
					placeholder="RENAVAM"
					value={renavam}
					onChange={(e) => setRenavam(e.target.value)}
					className="w-full px-4 py-3 rounded-lg border border-slate-300 focus:ring-2 focus:ring-[#1E3A8A]"
				/>
			</div>

			<VerifyButton plate={plate} renavam={renavam} wallet={wallet} onStart={onStart} />
		</div>
	);
}
