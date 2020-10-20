import compact from 'lodash/compact';
import fetch from 'node-fetch';
import fs from 'fs';
import request from 'request';
import { URLSearchParams } from 'url';
//import FormData from '@postman/form-data';
import path from 'path';

import {
  clientKeys,
  getApiHost,
  getAuthorizationEndpoint,
  namespace,
} from '../config';
import type { Envs } from '../types';

export async function getToken(env: Envs): Promise<string> {
  try {
    const authorizationEndpoint = getAuthorizationEndpoint();
    const clientId = clientKeys[env].clientId;
    const clientSecret = clientKeys[env].clientSecret;
    const body = new URLSearchParams();
    body.append('client_id', clientId);
    body.append('client_secret', clientSecret);
    body.append('grant_type', 'client_credentials');

    const res = await fetch(authorizationEndpoint, {
      method: 'POST',
      body,
    });

    if (res.ok) {
      const json = await res.json();
      return json.access_token;
    } else {
      throw new Error(res.statusText);
    }
  } catch (err) {
    throw new Error(err);
  }
}

type PublishParams = {
  configFile: string;
  dryRun: string;
  env: string;
  token: string;
};

type PublishResponse = {
  message: string;
  results: string;
};

// Temporarily using request due to an issue with FormData and save actions
export function publish({
  configFile,
  dryRun,
  env,
  token,
}: PublishParams): Promise<PublishResponse> {
  const apiHost = getApiHost(env);
  const filePath = path.resolve(process.cwd(), configFile);
  console.log(apiHost, token, namespace);
  const options = {
    method: 'PUT',
    url: `${apiHost}/namespaces/${namespace}/gateway`,
    headers: {
      Authorization: `Bearer ${token}`,
    },
    formData: {
      configFile: {
        value: fs.createReadStream(filePath),
        options: {
          filename: configFile,
          contentType: null,
        },
      },
      dryRun,
    },
  };

  return new Promise((resolve, reject) => {
    request(options, (error: Error, response: any) => {
      if (error) {
        reject(error);
      }
      const body = JSON.parse(response.body);

      if (response.statusCode >= 400) {
        reject({
          message: body.results || body.error,
        });
      }

      resolve(body);
    });
  });
}

type AddMembersParams = {
  env: Envs;
  users: string[] | undefined;
  managers: string[] | undefined;
};

export async function addMembers({
  env,
  users = [],
  managers = [],
}: AddMembersParams): Promise<any> {
  try {
    const apiHost = getApiHost(env);
    const token = await getToken(env);
    const usersToAdd = users.map((username: string) => ({
      username,
      roles: ['viewer'],
    }));
    const managersToAdd = managers.map((username: string) => ({
      username,
      roles: ['admin'],
    }));
    const body = compact([...managersToAdd, ...usersToAdd]);
    const res = await fetch(`${apiHost}/namespaces/${namespace}/membership`, {
      method: 'PUT',
      headers: {
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify(body),
    });
    const json = await res.json();

    if (res.ok) {
      return json;
    } else {
      throw new Error(json.error);
    }
  } catch (err) {
    throw new Error(err);
  }
}
