import type { Prompt } from '../../components/prompt-form';

const prompts: Prompt[] = [
  {
    label: 'OpenAPI JSON URL',
    key: 'url',
    constraint: {
      presence: { allowEmpty: false },
      url: true,
    },
  },
  {
    label: 'Route Host',
    key: 'routeHost',
    constraint: {
      presence: { allowEmpty: false },
    },
  },
  {
    label: 'Service URL',
    key: 'serviceUrl',
    constraint: {
      presence: { allowEmpty: false },
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

export default prompts;
