import LogoCarousel from './LogoCarousel';

export default function PoweredBy() {
	return (
		<section className="bg-white border-t border-slate-200">
			<div className="max-w-6xl mx-auto py-20 px-6 text-center">
				<p className="text-sm text-slate-500 font-semibold tracking-wide">POWERED BY</p>

				<div className="mt-12">
					<LogoCarousel />
				</div>
			</div>
		</section>
	);
}
