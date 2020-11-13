import { compile } from 'path-to-regexp';
import fetch, { RequestInit } from 'node-fetch';
import merge from 'lodash/merge';

import authenticate from './auth';
import {
  authorizationEndpoint,
  apiHost,
  env,
  legacyAuthorizationEndpoint,
  legacyApiHost,
  namespace,
} from '../config';

export function getEndpoints() {
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
    const config = merge(
      {},
      {
        method: 'GET',
        headers: {
          Authorization: `Bearer ${token}`,
        },
      },
      options
    );
    const res = await fetch(endpoint, config);

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
    const { auth, host } = getEndpoints();
    const token = await authenticate(auth);
    const uncompiledUrl = host + endpoint;
    let url = uncompiledUrl;

    if (endpoint.includes(':')) {
      const urlCompiler = compile(uncompiledUrl);
      url = urlCompiler({ namespace }, { validate: false });
    }
    const response = await api(token, url, options);
    return response;
  } catch (err) {
    console.log('hi', err);
    throw new Error(err);
  }
}

export default makeRequest;
