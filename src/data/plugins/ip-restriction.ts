import { IPlugin } from '../../types';
import { IIpRestriction } from './types';

const constraints = {
  allow: {
    type: 'array',
    presence: {
      allowEmpty: false,
    },
  },
};

const plugin: IPlugin = {
  id: 'ip-restriction',
  name: 'IP Restriction',
  description: 'lorem',
  enabled: false,
  constraints,
  config: {
    allow: ['10.10.10.0/24'],
  },
};

export default plugin;
