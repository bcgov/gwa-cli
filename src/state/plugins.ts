import create from 'zustand/vanilla';
import createHook from 'zustand';
import has from 'lodash/has';
import compact from 'lodash/compact';
import flow from 'lodash/flow';
import produce from 'immer';
import uniq from 'lodash/uniq';

import { loadPlugins } from '../services/plugins';
import type { PluginsResult } from '../types';

const store = create<PluginsResult>(() => ({}));

export const initPluginsState = async () => {
  try {
    const data = await loadPlugins();

    store.setState(produce((draft) => (draft = data)));
  } catch (err) {
    throw new Error(err);
  }
};

export const toggleEnabled = (id: string, enabled: boolean) =>
  store.setState(
    produce((draft) => {
      draft.data[id].meta.enabled = enabled;
    })
  );
export const set = (id: string, value: unknown) =>
  store.setState(
    produce((draft) => {
      draft.data[id].config = value;
    })
  );

export const generatePluginTemplates = (
  names: string[],
  namespace: string
): any[] => {
  const plugins = store.getState();
  const validate = flow(
    (value: string[]) => compact(value),
    (value: string[]) => uniq(value),
    (value: string[]) => value.filter((name: string) => has(plugins, name))
  );

  return validate(names).map((name) => ({
    name,
    tags: [`ns.${namespace}`],
    enabled: true,
    config: plugins[name].config,
  }));
};

export const usePluginsState = createHook(store);

export default store;
