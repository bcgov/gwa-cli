import type { Prompt } from '../../components/prompt-form';

const prompts: Prompt[] = [
  {
    label: 'Namespace',
    key: 'namespace',
    constraint: {
      presence: { allowEmpty: false },
      length: { minimum: 5, maximum: 10 },
      format: {
        pattern: '^[a-z][a-z0-9-]{4,14}$',
        flags: 'i',
        message: 'can only contain a-z and 0-9',
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
];

export default prompts;
