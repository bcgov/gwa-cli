export const namespace: string = process.env.GWA_NAMESPACE || '';
export const env = process.env.GWA_ENV || '';
export const clientId: string = process.env.CLIENT_ID || '';
export const clientSecret: string = process.env.CLIENT_SECRET || '';
export const authorizationEndpoint: string = `https://auth-qwzrwc-${env}.pathfinder.gov.bc.ca/auth/realms/aps/protocol/openid-connect/token`;
export const apiHost: string = `https://gwa-api-qwzrwc-${env}.pathfinder.gov.bc.ca/v1`;
