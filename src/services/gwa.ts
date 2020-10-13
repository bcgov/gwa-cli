import fetch from 'node-fetch';
import fs from 'fs';
import request from 'request';
import { URLSearchParams } from 'url';
import FormData from '@postman/form-data';
import path from 'path';

import {
  authorizationEndpoint,
  clientId,
  clientSecret,
  publishEndpoint,
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

type Publish = {
  configFile: string;
  dryRun: string;
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
  token,
}: Publish): Promise<PublishResponse> {
  const filePath = path.resolve(process.cwd(), configFile);
  const options = {
    method: 'PUT',
    url: publishEndpoint,
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
      if (error) reject(error);
      resolve(JSON.parse(response.body));
    });
  });
}
