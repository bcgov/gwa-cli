import { newRidgeState } from 'react-ridge-state';

export interface OrgState {
  name: string;
  specUrl: string;
  maintainers: string[];
}

export const orgState = newRidgeState<OrgState>({
  name: '',
  specUrl: '',
  maintainers: [],
});
