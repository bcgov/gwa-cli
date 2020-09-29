import path from 'path';
import dropRight from 'lodash/dropRight';

//import appState, { initAppState } from './state/app';
import { initPluginsState } from './state/plugins';
import { loadPlugins } from './services/plugins';

const run = async (fn: any, input: string, options: any) => {
  try {
    const data = await loadPlugins(path.resolve(__dirname, '../files'));
    initPluginsState(data);
    return fn(input, options);
  } catch (err) {
    console.error(err);
  }
};

export default run;
