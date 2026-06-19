import { Code, ConnectError } from "@connectrpc/connect";

// Base contract for application errors. Named AppError to avoid shadowing the
// global Error. `retryable()` lets callers decide whether to retry the operation.
export interface AppError {
  code: number;
  message: string;
  retryable(): boolean;
}

// Non-retryable: the request itself is wrong or conflicts with current state.
export class InvalidArgumentError implements AppError {
  readonly code = 400;
  constructor(readonly message = "入力内容が正しくありません") {}
  retryable(): boolean {
    return false;
  }
}

export class NotFoundError implements AppError {
  readonly code = 404;
  constructor(readonly message = "リソースが見つかりません") {}
  retryable(): boolean {
    return false;
  }
}

export class AlreadyExistError implements AppError {
  readonly code = 409;
  constructor(readonly message = "リソースは既に存在します") {}
  retryable(): boolean {
    return false;
  }
}

// Retryable: transient server / infrastructure failures.
export class InternalServerError implements AppError {
  readonly code = 500;
  constructor(readonly message = "サーバでエラーが発生しました") {}
  retryable(): boolean {
    return true;
  }
}

export class UnavailableError implements AppError {
  readonly code = 503;
  constructor(readonly message = "サービスが一時的に利用できません") {}
  retryable(): boolean {
    return true;
  }
}

// Translate any thrown value (typically a ConnectError) into an AppError.
export function toAppError(cause: unknown): AppError {
  const connectError = ConnectError.from(cause);
  switch (connectError.code) {
    case Code.InvalidArgument:
      return new InvalidArgumentError(connectError.message);
    case Code.NotFound:
      return new NotFoundError(connectError.message);
    case Code.AlreadyExists:
      return new AlreadyExistError(connectError.message);
    case Code.Unavailable:
      return new UnavailableError(connectError.message);
    default:
      return new InternalServerError(connectError.message);
  }
}
