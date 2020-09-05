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
  description:
    'Restrict access to a Service or a Route by either allowing or denying IP addresses. Single IPs, multiple IPs or ranges in CIDR notation like 10.10.10.0/24 can be used. The plugin supports IPv4 and IPv6 addresses.',
  constraints,
  data: {
    name: 'ip-restriction',
    enabled: false,
    config: {
      allow: ['10.10.10.0/24'],
    },
  },
};

export default plugin;
