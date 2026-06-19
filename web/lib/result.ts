import type { AppError } from "./errors";

// Application-wide result type. A discriminated union on `ok` so callers get
// `payload` narrowed on success and `error` narrowed on failure.
export type Result<T> =
  | { ok: true; payload: T }
  | { ok: false; error: AppError };

export function ok<T>(payload: T): Result<T> {
  return { ok: true, payload };
}

export function err<T>(error: AppError): Result<T> {
  return { ok: false, error };
}
