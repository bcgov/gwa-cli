export default {
  response_type: {
    type: 'string',
  },
  introspection_endpoint: {
    type: 'string',
    url: true,
  },
  filters: {
    type: 'string',
  },
  bearer_only: {
    type: 'string',
    inclusion: ['yes', 'no'],
  },
  ssl_verify: {
    type: 'string',
    inclusion: ['yes', 'no'],
  },
  session_secret: {
    type: 'string',
  },
  introspection_endpoint_auth_method: {
    type: 'string',
  },
  realm: {
    type: 'string',
  },
  redirect_after_logout_uri: {
    type: 'string',
  },
  scope: {
    type: 'string',
  },
  token_endpoint_auth_method: {
    type: 'string',
  },
  logout_path: {
    type: 'string',
  },
  client_id: {
    type: 'string',
    presence: {
      allowEmpty: false,
    },
  },
  client_secret: {
    type: 'string',
    presence: {
      allowEmpty: false,
    },
  },
  discovery: {
    type: 'string',
    url: true,
  },
  recovery_page_path: {
    type: 'string',
  },
  redirect_uri_path: {
    type: 'string',
  },
};
