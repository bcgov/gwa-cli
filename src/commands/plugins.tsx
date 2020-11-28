import { initPluginsState } from '../state/plugins';
import { loadPlugins } from '../services/plugins';
import render from '../views/plugins-list';

export default async function () {
  try {
    const data = await loadPlugins();
    initPluginsState(data);
    render();
  } catch (err) {
    throw new Error(err);
  }
}
