// src/app/page.tsx
import Navbar from '@/components/Navbar';
import TokenizationForm from '@/components/TokenizationForm';

export default function Home() {
	return (
		<>
			<Navbar />
			<div className="p-10">
				<TokenizationForm />
			</div>
		</>
	);
}
