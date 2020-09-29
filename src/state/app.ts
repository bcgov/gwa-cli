import create from 'zustand/vanilla';
import createStore from 'zustand';
import produce from 'immer';

export type AppState = {
  readonly cwd: string | null;
  readonly dir: string | null;
  input?: string | null;
  output?: string | null;
  mode: 'view' | 'edit';
  toggleMode: () => void;
};

const store = create<AppState>((set) => ({
  cwd: process.cwd(),
  dir: __dirname,
  input: null,
  output: null,
  mode: 'view',
  toggleMode: () =>
    set(
      produce((state) => {
        state.mode = state.mode === 'view' ? 'edit' : 'view';
      })
    ),
}));

export const initAppState = (
  data: Omit<AppState, 'toggleMode' | 'mode' | 'cwd' | 'dir'>
) => store.setState(data);

export const useAppState = createStore(store);

export default store;
