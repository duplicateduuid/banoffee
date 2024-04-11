import { JSDOM } from 'jsdom';
import { Readability } from '@mozilla/readability';
import { error } from '@sveltejs/kit';

export async function GET({ url }) {
	const urlParam = url.searchParams.get('url');

	if (!urlParam) {
		error(400, 'Missing url query param');
	}

	const doc = await JSDOM.fromURL(urlParam);
	const readability = new Readability(doc.window.document);
	const reader = readability.parse();

	return new Response(reader?.title || doc.window.document.title);
}
