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
    label: 'Output file name (.yaml)',
    key: 'outpfile',
    constraint: {
      presence: { allowEmpty: false },
      format: /^[\w,\s-]+\.(yaml|yml)/,
    },
  },
];

export const initialState = {
  step: 0,
  error: undefined,
  value: '',
  data: [],
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
        data: [...state.data, action.payload],
        step: state.step + 1,
        value: '',
      };

    case 'reset':
      return {
        ...state,
        value: '',
        error: undefined,
      };

    default:
      throw new Error();
  }
}

export default reducer;
