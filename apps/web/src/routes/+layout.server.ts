import type { LayoutServerLoad } from "./$types";
import { me } from "../requests/user";
import type { User } from "../schemas/user";

export const load: LayoutServerLoad<{ user: User | null}> = async ({ cookies }) => {  
  let user: User | null;
  
  try {
    user = await me(cookies);
  } catch (e) {
    console.error(e);
    user = null;
  }

  return { user };
} 
