import { createClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";
import { ProductService } from "@/gen/product/v1/product_pb";
import { type AppError, InternalServerError, toAppError } from "./errors";
import { sleep } from "./lib";
import { err, ok, type Result } from "./result";

// Server-side only: keep the API origin out of the browser bundle and avoid CORS
// by calling the gRPC server from Server Actions instead of the client.
const baseUrl = process.env.API_BASE_URL ?? "http://localhost:8080";

const transport = createConnectTransport({ baseUrl });

type RetryOption = {
  backOffTimeMs: number;
  maxRetry: number;
};

const defaultRetryOption: RetryOption = {
  backOffTimeMs: 300,
  maxRetry: 10,
};

// This is utility function wrapping gRPC client method.
// Client could be retryable, but I adopt this utility function from maintainability and readability.
export function withRetry<TArgs extends unknown[], TReturn>(
  fn: (...args: TArgs) => Promise<TReturn>,
  // Opt should be Partial, all fields are not mandatory, if one of them is undefined, default option is used.
  opt: Partial<RetryOption>,
): (...args: TArgs) => Promise<Result<TReturn>> {
  const retryOpt = {
    ...defaultRetryOption,
    ...opt,
  };
  return async (...args: TArgs) => {
    let lastErr: AppError = new InternalServerError();
    for (let i = 0; i < retryOpt.maxRetry; i++) {
      try {
        return ok(await fn(...args));
      } catch (e) {
        // Connect rejects with a thrown value (typically ConnectError), so the
        // failure path is a throw, not a resolved value. Normalize it first.
        lastErr = toAppError(e);
        // Retrying a non-retryable failure (e.g. NotFound) is pointless.
        if (!lastErr.retryable()) return err(lastErr);
        await sleep(retryOpt.backOffTimeMs);
      }
    }
    return err(lastErr);
  };
}
const rawProductClient = createClient(ProductService, transport);

// TODO: How to make this client testable?
// Mock function is a little bit hard to test with gRPC utility.
export const productClient = {
  getActiveProductsCount: withRetry(
    rawProductClient.getActiveProductsCount,
    {},
  ),
  listActiveProducts: withRetry(rawProductClient.listActiveProducts, {}),
  registerProduct: withRetry(rawProductClient.registerProduct, {}),
};
