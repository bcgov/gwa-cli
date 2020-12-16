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
          'https://authz-apps-gov-bc-ca.dev.apsgw.xyz/auth/realms/aps/protocol/openid-connect/token',
        apiHost: 'https://gwa-api-gov-bc-ca.dev.apsgw.xyz/v1',
      })
    );
  });

  it('should generate dynamic TEST urls', () => {
    process.env = { ...CACHED_ENV, GWA_ENV: 'test' };
    expect(config()).toEqual(
      expect.objectContaining({
        env: 'test',
        authorizationEndpoint:
          'https://authz-apps-gov-bc-ca.test.apsgw.xyz/auth/realms/aps/protocol/openid-connect/token',
        apiHost: 'https://gwa-api-gov-bc-ca.test.apsgw.xyz/v1',
      })
    );
  });

  it('should generate PROD urls', () => {
    process.env = { ...CACHED_ENV, GWA_ENV: 'prod' };
    expect(config()).toEqual(
      expect.objectContaining({
        env: 'prod',
        authorizationEndpoint: 'https://authz.apps.gov.bc.ca',
        apiHost: 'https://gwa.api.gov.bc.ca',
      })
    );
  });
});
