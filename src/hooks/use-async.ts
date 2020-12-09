import { useRef, useMemo } from 'react';
type State = 'pending' | 'error' | 'done';

export type PromiseFn<R, A extends any[] = []> = (...args: A) => Promise<R>;

export const cache = new Map();

const callPromise = <Response, Args extends any[]>(
  promise: PromiseFn<Response, Args>,
  ...args: Args
) => {
  const cached = cache.get(promise);

  if (cached) {
    return cached;
  }

  let result: Response;
  let error: any;
  let state: State = 'pending';

  const pending = promise(...args)
    .then((r: Response) => {
      result = r;
      state = 'done';
    })
    .catch((err) => {
      error = err;
      state = 'error';
    });

  const read = () => {
    if (state === 'pending') {
      throw pending;
    } else if (state === 'error') {
      throw new Error(error);
    }

    return result;
  };

  cache.set(promise, read);

  return read;
};

const useAsync = <Response>(
  promise: PromiseFn<Response, any[]>,
  ...args: any[]
) => {
  const result = useRef<Response | undefined>(undefined);

  useMemo(() => {
    result.current = callPromise(promise, ...args);
  }, [promise, ...args]);

  if (typeof result.current !== 'function') {
    return;
  }

  return result.current();
};

export default useAsync;
