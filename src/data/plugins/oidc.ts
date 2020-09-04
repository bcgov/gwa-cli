import { IPlugin } from '../../types';
import { IOidc } from './types';

const constraints = {
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
    presence: true,
  },
  client_secret: {
    type: 'string',
    presence: true,
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

const plugin: IPlugin = {
  id: 'oidc',
  name: 'OIDC',
  description: 'lorem',
  constraints,
  data: {
    name: 'oidc',
    enabled: false,
    config: {
      response_type: 'code',
      introspection_endpoint:
        'https://{HOST}/auth/realms/{REALM}/protocol/openid-connect/token/introspect',
      filters: null,
      bearer_only: 'no',
      ssl_verify: 'no',
      session_secret: null,
      introspection_endpoint_auth_method: null,
      realm: 'kong',
      redirect_after_logout_uri: '/',
      scope: 'openid',
      token_endpoint_auth_method: 'client_secret_post',
      logout_path: '/logout',
      client_id: '',
      client_secret: '',
      discovery:
        'https://{HOST}/auth/realms/{REALM}/.well-known/openid-configuration',
      recovery_page_path: null,
      redirect_uri_path: null,
    },
  },
};

export default plugin;
