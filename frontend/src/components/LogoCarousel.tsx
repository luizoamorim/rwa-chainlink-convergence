'use client';

const partners = [
	{ src: '/partners/tenderly.svg', size: 48 },
	{ src: '/partners/chainlink.svg', size: 48 },
	{ src: '/partners/worldid.svg', size: 152 },
	{ src: '/partners/thirdweb.svg', size: 32 },
];

export default function LogoCarousel() {
	const logos = [...partners, ...partners];

	return (
		<div className="relative overflow-hidden w-full mt-12">
			{/* fade left */}
			<div className="pointer-events-none absolute left-0 top-0 h-full w-32 bg-gradient-to-r from-white to-transparent z-10" />

			{/* fade right */}
			<div className="pointer-events-none absolute right-0 top-0 h-full w-32 bg-gradient-to-l from-white to-transparent z-10" />

			<div className="flex gap-24 animate-scroll items-center w-max">
				{logos.map((logo, i) => (
					<img key={i} src={logo.src} style={{ height: logo.size }} className="w-auto opacity-80 hover:opacity-100 transition" alt="partner" />
				))}
			</div>
		</div>
	);
}
