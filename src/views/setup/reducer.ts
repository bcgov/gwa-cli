const prompts = [
  {
    label: 'API Ownership Team',
    key: 'team',
    constraint: {
      presence: { allowEmpty: false },
    },
  },
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

export const initialState = {
  step: 0,
  error: undefined,
  status: 'idle',
  value: '',
  data: {},
  done: false,
  prompts,
};

function reducer(state, action) {
  switch (action.type) {
    case 'change':
      return {
        ...state,
        value: action.payload,
        error: undefined,
      };

    case 'error':
      return {
        ...state,
        error: action.payload,
      };

    case 'next':
      return {
        ...state,
        data: { ...state.data, ...action.payload },
        step: state.step + 1,
        value: '',
      };

    case 'reset':
      return {
        ...state,
        value: '',
        error: undefined,
      };
    case 'spec/loading':
      return {
        ...state,
        status: 'loading',
      };
    case 'spec/success':
      return {
        ...state,
        status: 'success',
      };

    case 'spec/failed':
      return {
        ...state,
        status: 'failed',
        specError: action.payload,
      };

    case 'spec/written':
      return {
        ...state,
        done: true,
      };

    default:
      throw new Error();
  }
}

export default reducer;
