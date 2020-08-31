import { IGwaEndpoint } from './types';
import { IPlugin } from '../../types';

const constraints = {
  api_owners: {
    type: 'array',
    presence: {
      allowEmpty: false,
    },
  },
};

const plugin: IPlugin = {
  id: 'bcgov-gwa-endpoint',
  name: 'bcgov-gwa-endpoint',
  enabled: false,
  constraints,
  config: {
    api_owners: [],
  },
};

export default plugin;
