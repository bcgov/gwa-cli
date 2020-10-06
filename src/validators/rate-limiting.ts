export default {
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
