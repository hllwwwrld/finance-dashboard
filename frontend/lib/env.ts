const apiBaseUrl = process.env.NEXT_PUBLIC_API_URL ?? "/api"

/** Next.js mock routes в frontend/app/api — только при `pnpm dev` без backend URL. */
export const isMockApiMode =
  process.env.NODE_ENV === "development" && !apiBaseUrl.startsWith("http")
