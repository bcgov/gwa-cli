import fetch from 'node-fetch';
import authenticate from '../auth';
import { api, makeRequest } from '../api';
import config from '../../config';

jest.mock('../auth');
jest.mock('node-fetch', () => require('fetch-mock-jest').sandbox());

const CACHED_ENV = process.env;

describe('services/api', () => {
  afterEach(() => {
    fetch.mockClear();
    fetch.mockReset();
  });

  describe('#api', () => {
    it('should call #authenticate first with endpoints', async () => {
      fetch.get('https://api.com/url', { test: 'object' });
      const response = await api('123', 'https://api.com/url');
      expect(response).toEqual({ test: 'object' });
    });

    it('should add auth header', async () => {
      fetch.get('https://api.com/url', { test: 'object' });
      const response = await api('123', 'https://api.com/url');
      expect(fetch.mock.calls[0][1]).toEqual({
        method: 'GET',
        headers: {
          Authorization: 'Bearer 123',
        },
      });
    });

    it('should take different methods', async () => {
      fetch.put('https://api.com/url', {});
      const response = await api('123', 'https://api.com/url', {
        method: 'PUT',
      });
      expect(fetch.mock.calls[0][1]).toEqual({
        method: 'PUT',
        headers: {
          Authorization: 'Bearer 123',
        },
      });
    });

    it('should throw error messages', async () => {
      fetch.get('https://api.com/url', 500);
      expect(async () => {
        await api('123', 'https://api.com/url');
      }).rejects.toThrow();
    });
  });

  describe('#makeRequest', () => {
    beforeEach(() => {
      process.env = { ...CACHED_ENV, GWA_ENV: 'dev', GWA_NAMESPACE: 'sampler' };
      jest.resetModules();
    });

    afterEach(() => {
      process.env = CACHED_ENV;
    });

    it('should call authenticate with the correct URL', async () => {
      const { apiHost, authorizationEndpoint } = config();
      fetch
        .post(authorizationEndpoint, { access_token: '123' })
        .get(`${apiHost}/test`, {});

      await makeRequest('/test');

      expect(authenticate).toHaveBeenCalledWith(authorizationEndpoint);
    });

    it('should parse parameters with env variables', async () => {
      const { apiHost, authorizationEndpoint } = config();
      fetch
        .post(authorizationEndpoint, { access_token: '123' })
        .get(`${apiHost}/sampler/endpoint`, {});

      await makeRequest('/:namespace/endpoint');

      expect(fetch.mock.calls[0][0]).toEqual(`${apiHost}/sampler/endpoint`);
    });

    it('throw an error', async () => {
      const { apiHost, authorizationEndpoint } = config();
      fetch
        .post(authorizationEndpoint, { access_token: '123' })
        .get(`${apiHost}/sampler/endpoint`, 500);

      expect(async () => {
        await makeRequest('/:namespace/endpoint');
      }).rejects.toThrow();
    });
  });
});
