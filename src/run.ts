import path from 'path';

import { initAppState } from './state/app';
import { initPluginsState } from './state/plugins';
import { loadPlugins } from './services/plugins';

const run = async (fn: any, input: string | null, options?: any) => {
  try {
    const data = await loadPlugins(path.resolve(__dirname, '../files'));
    initAppState({
      input,
      output: options?.output,
    });
    initPluginsState(data);

    return fn(input, options);
  } catch (err) {
    console.error(err);
  }
};

export default run;
