import type { Client } from "@connectrpc/connect";
import { Code, ConnectError } from "@connectrpc/connect";
import { beforeEach, describe, expect, it, vi } from "vitest";
import type { ProductService } from "@/gen/product/v1/product_pb";
import { makeRetryable } from "./client";

const rpc = vi.fn();

// The wrapped method is invoked as `wrapped.run(...)`; makeRetryable only cares
// that each property is a function returning a promise.
function wrap() {
  const client = { run: rpc } as unknown as Client<typeof ProductService>;
  return makeRetryable(client, { backOffInterval: 0 }) as unknown as {
    run: typeof rpc;
  };
}

describe("makeRetryable", () => {
  beforeEach(() => {
    rpc.mockReset();
  });

  it("returns ok with the payload on success and calls the rpc once", async () => {
    rpc.mockResolvedValue({ id: "p1" });

    const result = await wrap().run();

    expect(result).toEqual({ ok: true, payload: { id: "p1" } });
    expect(rpc).toHaveBeenCalledTimes(1);
  });

  it("does not retry a non-retryable error", async () => {
    rpc.mockRejectedValue(new ConnectError("dup", Code.AlreadyExists));

    const result = await wrap().run();

    expect(result.ok).toBe(false);
    if (!result.ok) {
      expect(result.error.retryable()).toBe(false);
    }
    expect(rpc).toHaveBeenCalledTimes(1);
  });

  it("retries a retryable error up to the max attempts then fails", async () => {
    rpc.mockRejectedValue(new ConnectError("down", Code.Unavailable));

    const result = await wrap().run();

    expect(result.ok).toBe(false);
    expect(rpc).toHaveBeenCalledTimes(3);
  });

  it("returns ok when a retry eventually succeeds", async () => {
    rpc
      .mockRejectedValueOnce(new ConnectError("down", Code.Unavailable))
      .mockResolvedValueOnce({ id: "p2" });

    const result = await wrap().run();

    expect(result).toEqual({ ok: true, payload: { id: "p2" } });
    expect(rpc).toHaveBeenCalledTimes(2);
  });
});
