import create from 'zustand/vanilla';
import createStore from 'zustand';
import produce from 'immer';

export type AppState = {
  readonly cwd: string | null;
  readonly dir: string | null;
  input?: string | null;
  output?: string | null;
  mode: 'view' | 'edit';
  team: string;
  toggleMode: () => void;
};

const store = create<AppState>((set) => ({
  cwd: process.cwd(),
  dir: __dirname,
  input: null,
  output: null,
  team: '',
  mode: 'view',
  toggleMode: () =>
    set(
      produce((state) => {
        state.mode = state.mode === 'view' ? 'edit' : 'view';
      })
    ),
}));

export const initAppState = (data: Pick<AppState, 'input' | 'output'>) =>
  store.setState(data);

export const useAppState = createStore(store);

export default store;
