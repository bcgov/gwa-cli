export const namespace: string = process.env.GWA_NAMESPACE || '';
export const env = process.env.GWA_ENV || '';
export const clientId: string = process.env.CLIENT_ID || '';
export const clientSecret: string = process.env.CLIENT_SECRET || '';
export const legacyAuthorizationEndpoint: string =
  'https://auth-qwzrwc-test.pathfinder.gov.bc.ca/auth/realms/aps/protocol/openid-connect/token';
export const legacyApiHost: string =
  'https://gwa-api-qwzrwc-test-env.pathfinder.gov.bc.ca/v1';
export const authorizationEndpoint: string = `https://auth-264e6f-${env}.apps.silver.devops.gov.bc.ca/auth/realms/aps/protocol/openid-connect/token`;
export const apiHost: string = `https://gwa-api-264e6f-${env}.apps.silver.devops.gov.bc.ca/v1`;
//TODO make this a function
