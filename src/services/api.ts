import { compile } from 'path-to-regexp';
import fetch, { RequestInit } from 'node-fetch';
import merge from 'lodash/merge';

import authenticate from './auth';
import config from '../config';

export function getEndpoints() {
  const {
    authorizationEndpoint,
    apiHost,
    env,
    legacyAuthorizationEndpoint,
    legacyApiHost,
  } = config();
  const isLegacy = env === 'test';
  const auth = isLegacy ? legacyAuthorizationEndpoint : authorizationEndpoint;
  const host = isLegacy ? legacyApiHost : apiHost;

  return {
    auth,
    host,
  };
}

export async function api<ApiResponse>(
  token: string,
  endpoint: string,
  options?: RequestInit
): Promise<ApiResponse> {
  try {
    const fetchOptions = merge(
      {},
      {
        method: 'GET',
        headers: {
          Authorization: `Bearer ${token}`,
        },
      },
      options
    );
    const res = await fetch(endpoint, fetchOptions);

    if (res.ok) {
      const json = await res.json();
      return json;
    } else {
      throw res.statusText;
    }
  } catch (err) {
    throw new Error(err);
  }
}

export async function makeRequest<ApiResponse>(
  endpoint: string,
  options?: RequestInit
): Promise<ApiResponse> {
  try {
    const { namespace } = config();
    const { auth, host } = getEndpoints();
    const token = await authenticate(auth);
    let path = endpoint;

    if (endpoint.includes(':')) {
      const compiler = compile(endpoint, {
        encode: encodeURIComponent,
      });
      path = compiler({ namespace });
    }
    const url = host + path;
    const response = await api<ApiResponse>(token, url, options);
    return response;
  } catch (err) {
    throw new Error(err);
  }
}

export default makeRequest;
