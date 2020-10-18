export const clientId: string = process.env.CLIENT_ID || '';
export const clientSecret: string = process.env.CLIENT_SECRET || '';
export const namespace: string = process.env.GWA_NAMESPACE || '';
export const service: string = process.env.GWA_SERVICE_NAME || '';

export const getAuthorizationEndpoint = (env: string = 'dev'): string =>
  `https://auth-qwzrwc-${env}.pathfinder.gov.bc.ca/auth/realms/aps/protocol/openid-connect/token`;
export const getApiHost = (env: string = 'dev'): string =>
  `https://gwa-api-qwzrwc-${env}.pathfinder.gov.bc.ca/v1`;
