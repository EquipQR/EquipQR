import { redirect } from "@sveltejs/kit"
import { getUserCurrent } from "$lib/api/auth"

export async function load({fetch}) {
    try {
        const user = await getUserCurrent(fetch);
        return { user };


    } catch (err) {
        throw redirect(302, "/portal/login");
    }
}