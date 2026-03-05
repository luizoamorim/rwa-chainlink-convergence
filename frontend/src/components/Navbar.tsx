'use client';

import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useEffect } from 'react';

import { ConnectButton, useActiveAccount } from 'thirdweb/react';
import { createThirdwebClient } from 'thirdweb';
import { createWallet } from 'thirdweb/wallets';

import { Car } from 'lucide-react';

const client = createThirdwebClient({
	clientId: process.env.NEXT_PUBLIC_THIRDWEB_CLIENT_ID!,
});

const wallets = [createWallet('io.metamask')];

export default function Navbar() {
	const account = useActiveAccount();
	const router = useRouter();

	// If user disconnects wallet, redirect to home page
	useEffect(() => {
		if (!account) {
			router.push('/');
		}
	}, [account, router]);

	return (
		<nav className="bg-white border-b border-slate-200 px-10 py-5 flex justify-between items-center">
			<Link href="/" className="flex items-center gap-3">
				<div className="bg-[#1E3A8A] p-2 rounded-xl">
					<Car className="text-white" size={22} />
				</div>

				<span className="text-xl font-semibold tracking-tight text-[#0B132B]">AutoLock DeFi</span>
			</Link>

			<div className="flex items-center gap-8">
				{/* Show only if connected */}

				{account && (
					<>
						<Link href="/tokenize" className="text-slate-700 hover:text-[#1E3A8A]">
							Tokenize
						</Link>

						<Link href="/assets" className="text-slate-700 hover:text-[#1E3A8A]">
							My Assets
						</Link>
					</>
				)}

				<ConnectButton client={client} wallets={wallets} theme="light" />
			</div>
		</nav>
	);
}
