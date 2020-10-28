import produce, { Draft } from 'immer';
import type { Prompt, SetupAction, SetupState } from './types';

export const makeInitialState = (prompts: Prompt[] = []): SetupState => ({
  step: 0,
  error: undefined,
  value: '',
  data: {},
  done: false,
  prompts,
});

const reducer = produce((draft: Draft<SetupState>, action: SetupAction) => {
  switch (action.type) {
    case 'change':
      draft.value = action.payload;
      draft.error = undefined;
      break;

    case 'error':
      draft.error = action.payload;
      break;

    case 'next':
      draft.data = { ...draft.data, ...action.payload };
      draft.step = draft.step + 1;
      draft.value = '';
      break;

    case 'done':
      draft.done = true;
      break;

    case 'reset':
      draft.value = '';
      draft.error = undefined;
      break;

    default:
      throw new Error();
  }
}, makeInitialState());

export default reducer;
