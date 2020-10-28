import compact from 'lodash/compact';
import fetch from 'node-fetch';
import fs from 'fs';
import request from 'request';
import { URLSearchParams } from 'url';
//import FormData from '@postman/form-data';
import { basename, extname, resolve } from 'path';

import {
  clientId,
  clientSecret,
  apiHost,
  authorizationEndpoint,
  namespace,
} from '../config';

const TEMP_FILE: string = '.temp.yaml';
const NAMESPACE_ERROR: string =
  'You do not have a namespace set. Check your .env file in this directory or run gwa init';

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
      throw res.statusText;
    }
  } catch (err) {
    throw new Error(err);
  }
}

export async function mergeConfigs() {
  const current = process.cwd();

  try {
    const dir = await fs.promises.opendir(current);
    const files = [];
    const separator = '\n---\n';

    for await (const dirent of dir) {
      const extension = extname(dirent.name);

      if (/\.(yaml|yml)/.test(extension)) {
        const file = await fs.promises.readFile(
          resolve(current, dirent.name),
          'utf8'
        );
        files.push(file);
      }
    }
    const result = await fs.promises.writeFile(
      TEMP_FILE,
      files.join(separator)
    );
    return result;
  } catch (err) {
    throw new Error(err);
  }
}

async function bundleFiles(configFile: string | undefined) {
  try {
    if (!configFile) {
      await mergeConfigs();
    }
    const filename = configFile || TEMP_FILE;
    const filePath = resolve(process.cwd(), filename);

    return fs.createReadStream(filePath);
  } catch (err) {
    throw new Error(err);
  }
}

type PublishParams = {
  configFile: string | undefined;
  dryRun: string;
  token: string;
};

type PublishResponse = {
  message: string;
  results: string;
};

// Temporarily using request due to an issue with FormData and save actions
export async function publish({
  configFile,
  dryRun,
  token,
}: PublishParams): Promise<PublishResponse> {
  const value = await bundleFiles(configFile);
  const options = {
    method: 'PUT',
    url: `${apiHost}/namespaces/${namespace}/gateway`,
    headers: {
      Authorization: `Bearer ${token}`,
    },
    formData: {
      configFile: {
        value,
        options: {
          filename: configFile,
          contentType: null,
        },
      },
      dryRun,
    },
  };

  return new Promise((resolve, reject) => {
    if (!namespace) {
      return reject(NAMESPACE_ERROR);
    }

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

      if (fs.existsSync(TEMP_FILE)) {
        fs.unlinkSync(TEMP_FILE);
      }

      resolve(body);
    });
  });
}

type AddMembersParams = {
  users: string[] | undefined;
  managers: string[] | undefined;
};

export async function addMembers({
  users = [],
  managers = [],
}: AddMembersParams): Promise<any> {
  try {
    if (!namespace) {
      throw NAMESPACE_ERROR;
    }

    const token = await getToken();
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
      throw json.error;
    }
  } catch (err) {
    throw new Error(err);
  }
}
