import compact from 'lodash/compact';
import fetch from 'node-fetch';
import fs from 'fs';
import request from 'request';
import { URLSearchParams } from 'url';
//import FormData from '@postman/form-data';
import path from 'path';

import {
  apiHost,
  authorizationEndpoint,
  clientId,
  clientSecret,
} from '../config';

export async function getToken(): Promise<string> {
  try {
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
  namespace: string;
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
  namespace,
  token,
}: PublishParams): Promise<PublishResponse> {
  const filePath = path.resolve(process.cwd(), configFile);
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
  namespace: string;
  users: string;
};

export async function addMembers({
  namespace,
  users,
}: AddMembersParams): Promise<any> {
  try {
    const token = await getToken();
    const body = compact(users.split(',')).map((user) => ({ username: user }));
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
