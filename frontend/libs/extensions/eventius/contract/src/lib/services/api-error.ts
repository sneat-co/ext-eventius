/** Pulls a human-readable message from an HttpErrorResponse carrying IApiError. */
export function extractErrorMessage(err: unknown, fallback: string): string {
  const body = (err as { error?: { message?: string } })?.error;
  return body?.message ?? fallback;
}
