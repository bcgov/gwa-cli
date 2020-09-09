import { newRidgeState } from 'react-ridge-state';

export interface OrgState {
  name: string;
  specUrl: string;
  host: string;
  file: string;
}

export const orgState = newRidgeState<OrgState>({
  name: '',
  specUrl: '',
  host: '',
  file: '',
});
