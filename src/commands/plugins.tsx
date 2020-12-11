import { initPluginsState } from '../state/plugins';
import render from '../views/plugins-list';

export default async function () {
  try {
    await initPluginsState();
    render();
  } catch (err) {
    throw new Error(err);
  }
}
