export const namespace: string = process.env.GWA_NAMESPACE || '';

type ClientKeyProps = {
  clientId: string;
  clientSecret: string;
};

export type ClientKeys = {
  dev: ClientKeyProps;
  test: ClientKeyProps;
  prod: ClientKeyProps;
};

export const clientKeys: ClientKeys = {
  dev: {
    clientId: process.env.DEV_CLIENT_ID || '',
    clientSecret: process.env.DEV_CLIENT_SECRET || '',
  },
  test: {
    clientId: process.env.TEST_CLIENT_ID || '',
    clientSecret: process.env.TEST_CLIENT_SECRET || '',
  },
  prod: {
    clientId: process.env.PROD_CLIENT_ID || '',
    clientSecret: process.env.PROD_CLIENT_SECRET || '',
  },
};

export const getAuthorizationEndpoint = (env: string = 'dev'): string =>
  `https://auth-qwzrwc-${env}.pathfinder.gov.bc.ca/auth/realms/aps/protocol/openid-connect/token`;
export const getApiHost = (env: string = 'dev'): string =>
  `https://gwa-api-qwzrwc-${env}.pathfinder.gov.bc.ca/v1`;
