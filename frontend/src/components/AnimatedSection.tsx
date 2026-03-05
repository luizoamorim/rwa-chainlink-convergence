'use client';

import { motion } from 'framer-motion';

export default function AnimatedSection({ children, bg = 'bg-white' }: { children: React.ReactNode; bg?: string }) {
	return (
		<section className={`${bg} py-32 border-t border-slate-200`}>
			<motion.div
				className="max-w-6xl mx-auto px-6"
				initial={{ opacity: 0, y: 60 }}
				whileInView={{ opacity: 1, y: 0 }}
				transition={{ duration: 0.7 }}
				viewport={{ once: false, amount: 0.3 }}
			>
				{children}
			</motion.div>
		</section>
	);
}
