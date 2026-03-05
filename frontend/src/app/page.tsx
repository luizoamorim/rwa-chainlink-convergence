import Navbar from '@/components/Navbar';
import Hero from '@/components/Hero';
import PoweredBy from '@/components/PoweredBy';
import HowItWorks from '@/components/HowItWorks';
import WhyAutoLock from '@/components/WhyAutoLock';
import Architecture from '@/components/Architecture';
import AnimatedSection from '@/components/AnimatedSection';

export default function Home() {
	return (
		<>
			<Hero />

			<AnimatedSection>
				<PoweredBy />
			</AnimatedSection>

			<AnimatedSection bg="bg-slate-50">
				<HowItWorks />
			</AnimatedSection>

			<AnimatedSection>
				<WhyAutoLock />
			</AnimatedSection>

			<AnimatedSection bg="bg-slate-50">
				<Architecture />
			</AnimatedSection>
		</>
	);
}
