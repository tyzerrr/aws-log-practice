import type { DescService } from "@bufbuild/protobuf";
import { type Client, createClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";
import { ProductService } from "@/gen/product/v1/product_pb";
import { type AppError, InternalServerError, toAppError } from "./errors";
import { err, ok, type Result } from "./result";

// Server-side only: keep the API origin out of the browser bundle and avoid CORS
// by calling the gRPC server from Server Actions instead of the client.
const baseUrl = process.env.API_BASE_URL ?? "http://localhost:8080";

const transport = createConnectTransport({ baseUrl });

export interface RetryOptions {
  maxAttempts: number;
  backOffInterval: number;
}

const defaultRetryOptions: RetryOptions = {
  maxAttempts: 3,
  backOffInterval: 100,
};

// Same surface as the underlying client, but every method retries retryable
// errors automatically and returns a Result instead of throwing.
export type RetryableClient<T extends DescService> = {
  [K in keyof Client<T>]: Client<T>[K] extends (
    ...args: infer A
  ) => Promise<infer R>
    ? (...args: A) => Promise<Result<R>>
    : Client<T>[K];
};

function sleep(ms: number): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

async function withRetry<R>(
  call: () => Promise<R>,
  options: RetryOptions,
): Promise<Result<R>> {
  let lastError: AppError = new InternalServerError();
  for (let attempt = 1; attempt <= options.maxAttempts; attempt++) {
    try {
      return ok(await call());
    } catch (cause) {
      lastError = toAppError(cause);
      // Only transient errors are worth retrying.
      if (!lastError.retryable()) {
        return err(lastError);
      }
      if (attempt < options.maxAttempts) {
        await sleep(options.backOffInterval * attempt);
      }
    }
  }
  return err(lastError);
}

// Wrap each RPC method so callers never write retry logic themselves: calling
// `client.someRpc(req)` retries based on the error type and yields a Result.
export function makeRetryable<T extends DescService>(
  client: Client<T>,
  options?: Partial<RetryOptions>,
): RetryableClient<T> {
  const resolved = { ...defaultRetryOptions, ...options };
  return new Proxy(client as object, {
    get(target, prop, receiver) {
      const original = Reflect.get(target, prop, receiver);
      if (typeof original !== "function") {
        return original;
      }
      return (...args: unknown[]) =>
        withRetry(
          () =>
            (original as (...a: unknown[]) => Promise<unknown>).apply(
              target,
              args,
            ),
          resolved,
        );
    },
  }) as unknown as RetryableClient<T>;
}

export const productClient = makeRetryable(
  createClient(ProductService, transport),
);
