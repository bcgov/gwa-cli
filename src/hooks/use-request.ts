type State = 'pending' | 'error' | 'done';

const makeRequest = <T>() => {
  let promise: Promise<T>;
  let state: State;
  let result: T;
  let error: Error | null;

  return (fn: any): T => {
    if (!promise) {
      promise = fn();
      state = 'pending';
      promise
        .then((json: T) => {
          state = 'done';
          result = json;
        })
        .catch((err: Error) => {
          state = 'error';
          error = err;
        });
    }

    if (state === 'pending') {
      throw promise;
    }

    if (state === 'error') {
      throw error;
    }

    if (state === 'done') {
      return result;
    }
    return result;
  };
};

export default makeRequest;
