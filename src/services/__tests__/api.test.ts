import fetch from 'node-fetch';
import { api, makeRequest } from '../api';
import config from '../../config';

jest.mock('node-fetch', () => require('fetch-mock-jest').sandbox());
jest.mock('../../config');

describe('services/api', () => {
  afterEach(() => {
    fetch.mockClear();
    fetch.mockReset();
  });

  describe('#api', () => {
    it('should call #authenticate first with endpoints', async () => {
      fetch.get('https://legacy-api.com/url', { test: 'object' });
      const response = await api('123', 'https://legacy-api.com/url');
      expect(response).toEqual({ test: 'object' });
    });

    it('should add auth header', async () => {
      fetch.get('https://legacy-api.com/url', { test: 'object' });
      const response = await api('123', 'https://legacy-api.com/url');
      expect(fetch.mock.calls[0][1]).toEqual({
        method: 'GET',
        headers: {
          Authorization: 'Bearer 123',
        },
      });
    });

    it('should take different methods', async () => {
      fetch.put('https://legacy-api.com/url', {});
      const response = await api('123', 'https://legacy-api.com/url', {
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
      fetch.get('https://legacy-api.com/url', 500);
      expect(async () => {
        await api('123', 'https://legacy-api.com/url');
      }).rejects.toThrow();
    });
  });

  describe('#makeRequest', () => {
    it('should call authenticate with the correct URL', async () => {
      fetch
        .post('https://auth.com/', { access_token: '123' })
        .get('https://api.com/test', {});
      const { makeRequest } = require('../api');

      await makeRequest('/test');

      expect(authenticate).toHaveBeenCalledWith('https://auth.com');
    });
  });

  describe('#getEndpoints', () => {
    it('should return legacy endpoints', () => {
      const { getEndpoints } = require('../api');
      expect(getEndpoints()).toEqual({
        auth: 'https://legacy-auth.com',
        host: 'https://legacy-api.com',
      });
    });

    it('should return current endpoints', () => {
      config.mockImplementation(() => ({
        env: 'dev',
      }));
      const { getEndpoints } = require('../api');

      expect(getEndpoints()).toEqual({
        auth: 'https://auth.com',
        host: 'https://api.com',
      });
    });
  });
});
