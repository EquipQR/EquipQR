import { redirect } from "@sveltejs/kit";
import { getUserCurrent } from "$lib/api/auth";

export async function load({ fetch }) {
  try {
    await getUserCurrent(fetch);
  } catch {
    return;
  }

  throw redirect(302, "/");
}
