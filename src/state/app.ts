import { newRidgeState } from 'react-ridge-state';

export interface AppState {
  mode: 'view' | 'edit';
}

export const appState = newRidgeState<AppState>({
  mode: 'view',
});
