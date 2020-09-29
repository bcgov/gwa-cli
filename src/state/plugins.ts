import create from 'zustand/vanilla';
import produce from 'immer';

import type { PluginsResult } from '../types';

const store = create<PluginsResult>(() => ({}));

export const initPluginsState = (data: PluginsResult) => store.setState(data);

export const set = (id: string, key: string, value: unknown) =>
  store.setState(
    produce((draft) => {
      draft.data[id].config[key] = value;
    })
  );

export default store;
