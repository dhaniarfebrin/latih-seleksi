// Message of a caught value; anything that is not an Error is stringified.
export const errorMessage = (e: unknown): string => (e instanceof Error ? e.message : String(e))
