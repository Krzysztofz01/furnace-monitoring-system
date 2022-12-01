export type Result =
    | { isSuccess: true }
    | { isSuccess: false };

export type ValueResult<TResult, TError = Error> =
    | { isSuccess: true, value: TResult }
    | { isSuccess: false; error: TError };