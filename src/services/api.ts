import { compile } from 'path-to-regexp';
import fetch, { RequestInit } from 'node-fetch';
import merge from 'lodash/merge';

import authenticate from './auth';
import config from '../config';

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
    const { authorizationEndpoint, apiHost, namespace } = config();
    const token = await authenticate(authorizationEndpoint);
    let path = endpoint;

    if (endpoint.includes(':')) {
      const compiler = compile(endpoint, {
        encode: encodeURIComponent,
      });
      path = compiler({ namespace });
    }
    const url = apiHost + path;
    const response = await api<ApiResponse>(token, url, options);
    return response;
  } catch (err) {
    throw new Error(err);
  }
}

export default makeRequest;
