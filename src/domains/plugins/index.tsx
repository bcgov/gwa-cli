//import { render } from 'ink';

import { load } from './state';
import { loadPlugins } from '../../services/plugins';

export default async function (path: string) {
  try {
    const data = await loadPlugins(path);
    load(data);
    console.log(data);
  } catch (err) {
    console.error(err);
  }
}
