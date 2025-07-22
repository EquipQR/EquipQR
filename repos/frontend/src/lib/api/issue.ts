import type { CreateIssueRequest } from "$lib/types/issue";

export async function createIssue(
	request: CreateIssueRequest,
	fetchFn: typeof fetch = fetch
): Promise<Response> {
	const res = await fetchFn('/api/issue', {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json',
		},
		body: JSON.stringify(request),
	});

	if (!res.ok) {
		const err = await res.json();
		throw new Error(err.error || 'Failed to create issue');
	}

	return res;
}
