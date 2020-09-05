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
  name: 'GWA Endpoint',
  description:
    'Kong Plugin to process BC Government siteminder headers to apply kong consumers and acls (groups) to BC Government users.',
  constraints,
  data: {
    name: 'bcgov-gwa-endpoint',
    enabled: false,
    config: {
      api_owners: [],
    },
  },
};

export default plugin;
