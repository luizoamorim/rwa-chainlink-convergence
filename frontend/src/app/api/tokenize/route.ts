import { NextRequest, NextResponse } from 'next/server';

export async function POST(req: NextRequest) {
	try {
		const body = await req.json();

		const response = await fetch('http://localhost:8081/tokenize', {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify(body),
		});

		const text = await response.text();

		if (!response.ok) {
			return NextResponse.json({ error: text }, { status: 500 });
		}

		return NextResponse.json({ success: true });
	} catch (err: any) {
		return NextResponse.json({ error: err.message }, { status: 500 });
	}
}
