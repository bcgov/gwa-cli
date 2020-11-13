import fetch from 'node-fetch';

import authenticate from '../auth';

jest.mock('node-fetch', () => require('fetch-mock-jest').sandbox());

describe('services/auth', () => {
  const AUTH_URL = 'http://url.com';

  afterEach(() => {
    fetch.mockClear();
    fetch.mockReset();
  });

  describe('#authenticate', () => {
    it('should fetch with correct options', async () => {
      fetch.post(AUTH_URL, { access_token: '123' });
      const token = await authenticate(AUTH_URL);

      expect(fetch).toHaveBeenCalledWith(
        AUTH_URL,
        expect.objectContaining({
          body: expect.anything(),
          method: 'POST',
        })
      );
      expect(token).toEqual('123');
    });

    it('should handle failures', () => {
      fetch.post(AUTH_URL, 500);
      expect.assertions(1);
      expect(async () => {
        await authenticate(AUTH_URL);
      }).rejects.toThrow();
    });
  });
});
