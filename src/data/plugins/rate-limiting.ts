import { IPlugin } from '../../types';
import { IRateLimiting } from './types';

const constraints = {
  fault_tolerant: {
    type: 'boolean',
  },
  hide_client_headers: {
    type: 'boolean',
  },
  limit_by: {
    type: 'string',
    inclusion: ['consumer', 'credential', 'ip', 'service', 'header'],
  },
  minute: {
    type: 'number',
  },
  policy: {
    type: 'string',
  },
  header_name: {
    type: 'string',
  },
  second: {
    type: 'number',
  },
  hour: {
    type: 'number',
  },
  day: {
    type: 'number',
  },
  month: {
    type: 'number',
  },
  year: {
    type: 'number',
  },
  redis_database: {
    type: 'number',
  },
  redis_host: {
    type: 'string',
  },
  redis_password: {
    type: 'number',
  },
  redis_port: {
    type: 'number',
  },
  redis_timeout: {
    type: 'number',
  },
};

const plugin: IPlugin = {
  id: 'rate-limiting',
  name: 'Rate Limiting',
  description:
    'Rate limit how many HTTP requests can be made in a given period of seconds, minutes, hours, days, months, or years. If the underlying Service/Route (or deprecated API entity) has no authentication layer, the Client IP address will be used, otherwise the Consumer will be used if an authentication plugin has been configured.',
  constraints,
  encrypted: [],
  data: {
    name: 'rate-limiting',
    enabled: false,
    config: {
      fault_tolerant: true,
      hide_client_headers: false,
      limit_by: 'consumer',
      minute: 10,
      policy: 'cluster',
      header_name: null,
      second: null,
      hour: null,
      day: null,
      month: null,
      year: null,
      redis_database: 0,
      redis_host: null,
      redis_password: null,
      redis_port: 6379,
      redis_timeout: 2000,
    },
  },
};

export default plugin;
