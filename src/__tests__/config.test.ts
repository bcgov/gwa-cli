import config from '../config';

const CACHED_ENV = process.env;

describe('config', () => {
  beforeAll(() => {
    process.env = { ...CACHED_ENV, GWA_NAMESPACE: 'sampler' };
  });

  afterAll(() => {
    process.env = CACHED_ENV;
    jest.resetAllMocks();
  });

  it('should generate dynamic DEV urls', () => {
    process.env = { ...CACHED_ENV, GWA_ENV: 'dev' };
    expect(config()).toEqual(
      expect.objectContaining({
        env: 'dev',
        authorizationEndpoint:
          'https://authz-apps-gov-bc-ca.dev.api.gov.bc.ca/auth/realms/aps/protocol/openid-connect/token',
        apiHost: 'https://gwa-api-gov-bc-ca.dev.api.gov.bc.ca/v2',
        apiVersion: '2',
      })
    );
  });

  it('should generate dynamic TEST urls', () => {
    process.env = { ...CACHED_ENV, GWA_ENV: 'test' };
    expect(config()).toEqual(
      expect.objectContaining({
        env: 'test',
        authorizationEndpoint:
          'https://authz-apps-gov-bc-ca.test.api.gov.bc.ca/auth/realms/aps/protocol/openid-connect/token',
        apiHost: 'https://gwa-api-gov-bc-ca.test.api.gov.bc.ca/v2',
        apiVersion: '2',
      })
    );
  });

  it('should generate PROD urls', () => {
    process.env = { ...CACHED_ENV, GWA_ENV: 'prod' };
    expect(config()).toEqual(
      expect.objectContaining({
        env: 'prod',
        authorizationEndpoint:
          'https://authz.apps.gov.bc.ca/auth/realms/aps/protocol/openid-connect/token',
        apiHost: 'https://gwa.api.gov.bc.ca/v2',
        apiVersion: '2',
      })
    );
  });
});
