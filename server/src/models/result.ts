export type Result =
    | { isSuccess: true }
    | { isSuccess: false };

export type ValueResult<TResult> =
    | { isSuccess: true, value: TResult }
    | { isSuccess: false; value: undefined };