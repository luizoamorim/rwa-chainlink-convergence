// src/app/layout.tsx

import './globals.css';
import { ReactNode } from 'react';
import Providers from '@/components/Providers';
import Navbar from '@/components/Navbar';

export default function RootLayout({ children }: { children: ReactNode }) {
	return (
		<html lang="en">
			<body>
				<Providers>
					<Navbar />
					{children}
				</Providers>
			</body>
		</html>
	);
}
