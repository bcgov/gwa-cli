import path from 'path';

import { initAppState } from './state/app';

const run = async (fn: any, input: string | null, options?: any) => {
  try {
    initAppState({
      input,
      output: options?.output,
    });

    return fn(input, options);
  } catch (err) {
    console.error(err);
  }
};

export default run;
