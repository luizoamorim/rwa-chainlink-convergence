'use client';

import { ConnectButton } from 'thirdweb/react';
import { createThirdwebClient } from 'thirdweb';
import { createWallet } from 'thirdweb/wallets';
import { Car } from 'lucide-react';

const client = createThirdwebClient({
	clientId: process.env.NEXT_PUBLIC_THIRDWEB_CLIENT_ID!,
});

const wallets = [createWallet('io.metamask')];

export default function Navbar() {
	return (
		<nav className="bg-white border-b border-slate-200 px-10 py-5 flex justify-between items-center shadow-sm">
			<div className="flex items-center gap-3">
				<div className="bg-blue-600 p-2 rounded-xl shadow-md">
					<Car className="text-white" size={22} />
				</div>
				<span className="text-xl font-semibold tracking-tight text-slate-900">AutoLock DeFi</span>
			</div>

			<ConnectButton client={client} wallets={wallets} theme="light" />
		</nav>
	);
}
