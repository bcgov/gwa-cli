import create from 'zustand/vanilla';
import createHook from 'zustand';
import produce from 'immer';

import type { PluginsResult } from '../types';

const store = create<PluginsResult>(() => ({}));

export const initPluginsState = (data: PluginsResult) => store.setState(data);

export const loadPlugins = (data: any[]) => {
  store.setState(
    produce((draft) => {
      data.forEach((plugin) => {
        draft[plugin.name].config = plugin.config;
      });
    })
  );
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

export const usePluginsState = createHook(store);

export default store;
