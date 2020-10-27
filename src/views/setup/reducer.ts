import produce, { Draft } from 'immer';
import type { Prompt, SetupAction, SetupState } from './types';

const prompts: Prompt[] = [
  {
    label: 'OpenAPI JSON URL',
    key: 'url',
    constraint: {
      presence: true,
      url: true,
    },
  },
  {
    label: 'Starter plugins',
    key: 'plugins',
    constraint: {},
  },
  {
    label: 'Output file name (.yaml)',
    key: 'outfile',
    constraint: {
      presence: { allowEmpty: false },
      format: /^[\w,\s-]+\.(yaml|yml)/,
    },
  },
];

export const initialState: SetupState = {
  step: 0,
  error: undefined,
  status: 'idle',
  value: '',
  data: {},
  done: false,
  prompts,
  specError: '',
};

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

    case 'reset':
      draft.value = '';
      draft.error = undefined;
      break;

    case 'spec/loading':
      draft.status = 'loading';
      break;

    case 'spec/success':
      draft.status = 'success';
      break;

    case 'spec/failed':
      draft.status = 'failed';
      draft.specError = action.payload;
      break;

    case 'spec/written':
      draft.done = true;
      break;

    default:
      throw new Error();
  }
}, initialState);

export default reducer;
