import { newRidgeState } from 'react-ridge-state';

export const orgState = newRidgeState<any>({
  name: '',
  specUrl: '',
  maintainers: [],
});
