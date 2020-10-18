import path from 'path';
import pick from 'lodash/pick';

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
    const envFlags = pick(options, ['dev', 'prod', 'test']);
    const keys = Object.keys(envFlags);

    return fn(input, { ...options, env: keys[0] || 'dev' });
  } catch (err) {
    console.error(err);
  }
};

export default run;
