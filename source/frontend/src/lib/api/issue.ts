import type { CreateIssueRequest } from "$lib/types/issue";

export async function createIssueAndUploadFiles(
  request: CreateIssueRequest,
  fetchFn: typeof fetch = fetch
): Promise<Response> {
  // Step 1: Create the issue
  const res = await fetchFn("/api/issue", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      title: request.title,
      description: request.description,
      equipmentId: request.equipmentId,
    }),
  });

  let issueData: { id?: string };

  try {
    issueData = await res.json();
  } catch {
    throw new Error("Failed to parse issue creation response");
  }

  if (!res.ok) {
    console.error("Issue creation failed:", issueData);
    throw new Error(issueData?.error ?? "Failed to create issue");
  }

  const issueId = issueData.id;
  if (!issueId) {
    console.error("Missing issue_id in response:", issueData);
    throw new Error("Issue ID missing in response");
  }

  // Step 2: Upload attachments if present
  if (request.files?.length) {
    const formData = new FormData();
    for (const file of request.files) {
      formData.append("files", file);
    }

    const uploadRes = await fetchFn(`/api/issue/${issueId}/attachments`, {
      method: "POST",
      body: formData,
    });

    if (!uploadRes.ok) {
      const err = await uploadRes.json().catch(() => ({}));
      console.error("File upload failed:", err);
      throw new Error(err?.error || "Failed to upload attachments");
    }
  }

  return res;
}
