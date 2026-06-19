import { Code, ConnectError } from "@connectrpc/connect";
import { describe, expect, it } from "vitest";
import {
  AlreadyExistError,
  InternalServerError,
  InvalidArgumentError,
  NotFoundError,
  toAppError,
  UnavailableError,
} from "./errors";

describe("application errors", () => {
  it("marks client/state errors as non-retryable", () => {
    expect(new AlreadyExistError().retryable()).toBe(false);
    expect(new NotFoundError().retryable()).toBe(false);
    expect(new InvalidArgumentError().retryable()).toBe(false);
  });

  it("marks transient server errors as retryable", () => {
    expect(new InternalServerError().retryable()).toBe(true);
    expect(new UnavailableError().retryable()).toBe(true);
  });
});

describe("toAppError", () => {
  it("maps Code.AlreadyExists to a non-retryable AlreadyExistError", () => {
    const error = toAppError(new ConnectError("dup", Code.AlreadyExists));

    expect(error).toBeInstanceOf(AlreadyExistError);
    expect(error.retryable()).toBe(false);
  });

  it("maps Code.NotFound to a non-retryable NotFoundError", () => {
    const error = toAppError(new ConnectError("missing", Code.NotFound));

    expect(error).toBeInstanceOf(NotFoundError);
    expect(error.retryable()).toBe(false);
  });

  it("maps Code.Unavailable to a retryable UnavailableError", () => {
    const error = toAppError(new ConnectError("down", Code.Unavailable));

    expect(error).toBeInstanceOf(UnavailableError);
    expect(error.retryable()).toBe(true);
  });

  it("maps unknown causes to a retryable InternalServerError", () => {
    const error = toAppError(new Error("boom"));

    expect(error).toBeInstanceOf(InternalServerError);
    expect(error.retryable()).toBe(true);
  });
});
