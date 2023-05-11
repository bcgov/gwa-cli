import type { Prompt } from '../../components/prompt-form';

const prompts: Prompt[] = [
  {
    label: 'Namespace',
    key: 'namespace',
    constraint: {
      presence: { allowEmpty: false },
      length: { minimum: 5, maximum: 15 },
      format: {
        pattern: '^[a-z][a-z0-9-]{4,14}$',
        flags: 'gi',
        message: 'can only contain a-z, 0-9 and dashes',
      },
    },
  },
  {
    label: 'Client ID',
    key: 'clientId',
    constraint: {
      presence: { allowEmpty: false },
    },
  },
  {
    label: 'Client Secret',
    key: 'clientSecret',
    secret: true,
    constraint: {
      presence: { allowEmpty: false },
    },
  },
  {
    label: 'API Version',
    key: 'apiVersion',
    constraint: {
      format: {
        pattern: '[1-2]+',
        message: 'only versions 1 and 2 are available',
      },
    },
  },
  {
    label: 'Directory API Version',
    key: 'dirApiVersion',
    constraint: {
      format: {
        pattern: '[1-2]+',
        message: 'only versions 1 and 2 are available',
      },
    },
  },
  {
    label: 'Data Center (optional)',
    key: 'dataCenter',
    constraint: {
      presence: { allowEmpty: true },
    },
  },
];

export default prompts;
