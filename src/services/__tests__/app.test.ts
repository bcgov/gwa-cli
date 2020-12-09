import mock from 'mock-fs';

import * as app from '../app';

describe('services/app', () => {
  afterEach(() => {
    mock.restore();
  });

  describe('#checkForEnvFile', () => {
    it('should return false for missing', () => {
      expect(app.checkForEnvFile()).toEqual(false);
    });

    it('should return false for missing', () => {
      mock({
        '.env': 'content',
      });
      expect(app.checkForEnvFile()).toEqual(true);
    });
  });

  describe('#makeEnvFile', () => {
    it('should write an env file', () => {
      expect(
        app.makeEnvFile({
          namespace: 'sampler',
          clientId: 'id',
          clientSecret: 'secret',
          env: 'dev',
        })
      ).resolves.toEqual('.env file successfully generated');
    });
  });
});
