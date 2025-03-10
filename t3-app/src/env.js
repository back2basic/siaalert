import { createEnv } from "@t3-oss/env-nextjs";
import { z } from "zod";

export const env = createEnv({
  /**
   * Specify your server-side environment variables schema here. This way you can ensure the app
   * isn't built with invalid env vars.
   */
  server: {
    NODE_ENV: z.enum(["development", "test", "production"]),
  },

  /**
   * Specify your client-side environment variables schema here. This way you can ensure the app
   * isn't built with invalid env vars. To expose them to the client, prefix them with
   * `NEXT_PUBLIC_`.
   */
  client: {
    // NEXT_PUBLIC_CLIENTVAR: z.string(),
    // NEXT_PUBLIC_APPWRITE_ENDPOINT: z.string().url(),
    // NEXT_PUBLIC_APPWRITE_PROJECTID: z.string(),

    // NEXT_PUBLIC_APPWRITE_DBZEN: z.string(),
    // NEXT_PUBLIC_APPWRITE_COLHOSTS_ZEN: z.string(),
    // NEXT_PUBLIC_APPWRITE_COLSTATUS_ZEN: z.string(),
    // NEXT_PUBLIC_APPWRITE_COLCHECK_ZEN: z.string(),

    // NEXT_PUBLIC_APPWRITE_DBMAIN: z.string(),
    // NEXT_PUBLIC_APPWRITE_COLHOSTS_MAIN: z.string(),
    // NEXT_PUBLIC_APPWRITE_COLSTATUS_MAIN: z.string(),
    // NEXT_PUBLIC_APPWRITE_COLCHECK_MAIN: z.string(),
    NEXT_PUBLIC_NETWORK_ZEN_URL: z.string().url(),
    NEXT_PUBLIC_NETWORK_MAIN_URL: z.string().url(),
  },

  /**
   * You can't destruct `process.env` as a regular object in the Next.js edge runtimes (e.g.
   * middlewares) or client-side so we need to destruct manually.
   */
  runtimeEnv: {
    NODE_ENV: process.env.NODE_ENV,
    // NEXT_PUBLIC_CLIENTVAR: process.env.NEXT_PUBLIC_CLIENTVAR,
    // NEXT_PUBLIC_APPWRITE_ENDPOINT: process.env.NEXT_PUBLIC_APPWRITE_ENDPOINT,
    // NEXT_PUBLIC_APPWRITE_PROJECTID: process.env.NEXT_PUBLIC_APPWRITE_PROJECTID,
    // NEXT_PUBLIC_APPWRITE_DBZEN: process.env.NEXT_PUBLIC_APPWRITE_DBZEN,
    // NEXT_PUBLIC_APPWRITE_COLHOSTS_ZEN:
    //   process.env.NEXT_PUBLIC_APPWRITE_COLHOSTS_ZEN,
    // NEXT_PUBLIC_APPWRITE_COLSTATUS_ZEN:
    //   process.env.NEXT_PUBLIC_APPWRITE_COLSTATUS_ZEN,
    // NEXT_PUBLIC_APPWRITE_COLCHECK_ZEN:
    //   process.env.NEXT_PUBLIC_APPWRITE_COLCHECK_ZEN,
    // NEXT_PUBLIC_APPWRITE_DBMAIN: process.env.NEXT_PUBLIC_APPWRITE_DBMAIN,
    // NEXT_PUBLIC_APPWRITE_COLHOSTS_MAIN:
    //   process.env.NEXT_PUBLIC_APPWRITE_COLHOSTS_MAIN,
    // NEXT_PUBLIC_APPWRITE_COLSTATUS_MAIN:
    //   process.env.NEXT_PUBLIC_APPWRITE_COLSTATUS_MAIN,
    // NEXT_PUBLIC_APPWRITE_COLCHECK_MAIN:
    //   process.env.NEXT_PUBLIC_APPWRITE_COLCHECK_MAIN,
    NEXT_PUBLIC_NETWORK_ZEN_URL: process.env.NEXT_PUBLIC_NETWORK_ZEN_URL,
    NEXT_PUBLIC_NETWORK_MAIN_URL: process.env.NEXT_PUBLIC_NETWORK_MAIN_URL,
  },
  /**
   * Run `build` or `dev` with `SKIP_ENV_VALIDATION` to skip env validation. This is especially
   * useful for Docker builds.
   */
  skipValidation: !!process.env.SKIP_ENV_VALIDATION,
  /**
   * Makes it so that empty strings are treated as undefined. `SOME_VAR: z.string()` and
   * `SOME_VAR=''` will throw an error.
   */
  emptyStringAsUndefined: true,
});
