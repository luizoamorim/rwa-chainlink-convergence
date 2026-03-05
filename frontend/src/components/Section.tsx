export default function Section({ children, bg = 'bg-white' }: { children: React.ReactNode; bg?: string }) {
	return (
		<section className={`${bg} min-h-[70vh] flex items-center border-t border-slate-200`}>
			<div className="max-w-6xl mx-auto px-6 w-full">{children}</div>
		</section>
	);
}
