import create from 'zustand/vanilla';
import produce from 'immer';

export type AppState = {
  readonly cwd: string | null;
  readonly dir: string | null;
  input?: string | null;
  output?: string | null;
  mode: 'view' | 'edit';
};

const store = create<AppState>(() => ({
  cwd: process.cwd(),
  dir: __dirname,
  input: null,
  output: null,
  mode: 'view',
}));

export const initAppState = (data: Omit<AppState, 'mode' | 'cwd' | 'dir'>) =>
  store.setState(data);

export default store;
