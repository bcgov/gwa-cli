import { IGwaIpAnonymity } from './types';
import { IPlugin } from '../../types';

const constraints = {
  ipv4_mask: {
    type: 'number',
    presence: true,
  },
  ipv6_mask: {
    type: 'number',
    presence: true,
  },
};

const plugin: IPlugin = {
  id: 'gwa-ip-anonymity',
  name: 'IP Anonymity',
  description: 'lorem',
  constraints,
  data: {
    name: 'gwa-ip-anonymity',
    enabled: false,
    config: {
      ipv4_mask: 0,
      ipv6_mask: 0,
    },
  },
};

export default plugin;
