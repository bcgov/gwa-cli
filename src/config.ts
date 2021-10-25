function config() {
  const namespace: string = process.env.GWA_NAMESPACE || '';
  const env: string = process.env.GWA_ENV || '';
  const dataCenter: string = process.env.DATA_CENTER || '';
  const clientId: string = process.env.CLIENT_ID || '';
  const clientSecret: string = process.env.CLIENT_SECRET || '';
  const apiVersion: string = process.env.API_VERSION || '2';
  const legacyAuthorizationEndpoint: string =
    'https://auth-qwzrwc-test.pathfinder.gov.bc.ca/auth/realms/aps/protocol/openid-connect/token';
  const legacyApiHost: string = `https://gwa-api-qwzrwc-test.pathfinder.gov.bc.ca/v${apiVersion}`;
  let authorizationEndpoint: string = `https://authz-apps-gov-bc-ca.${env}.api.gov.bc.ca/auth/realms/aps/protocol/openid-connect/token`;
  let apiHost: string = `https://gwa-api-gov-bc-ca.${env}.api.gov.bc.ca/v${apiVersion}`;
  const dsApiHost: string = `https://api-gov-bc-ca.${env}.api.gov.bc.ca`;

  if (env === 'prod') {
    authorizationEndpoint =
      'https://authz.apps.gov.bc.ca/auth/realms/aps/protocol/openid-connect/token';
    apiHost = `https://gwa.api.gov.bc.ca/v${apiVersion}`;
  }

  if (dataCenter === 'kdc' || dataCenter === 'cdc') {
    authorizationEndpoint = `https://authz-apps-gov-bc-ca-${env}.${dataCenter}.api.gov.bc.ca/auth/realms/aps/protocol/openid-connect/token`;
    apiHost = `https://gwa-api-gov-bc-ca-${env}.${dataCenter}.api.gov.bc.ca/v${apiVersion}`;
    dsApiHost = `https://api-gov-bc-ca-${env}.${dataCenter}.api.gov.bc.ca`;
  }

  return {
    apiVersion,
    authorizationEndpoint,
    apiHost,
    clientId,
    clientSecret,
    dsApiHost,
    env,
    namespace,
    legacyAuthorizationEndpoint,
    legacyApiHost,
  };
}

export type Config = ReturnType<typeof config>;

export default config;
