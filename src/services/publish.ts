import { compile } from 'path-to-regexp';
import fs from 'fs';
import request from 'request';
import { extname, resolve } from 'path';
import YAML from 'yaml';

import authenticate from './auth';
import config from '../config';

const { namespace } = config();
const TEMP_FILE: string = '.temp.yaml';
const NAMESPACE_ERROR: string =
  'You do not have a namespace set. Check your .env file in this directory or run gwa init';

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

        try {
          YAML.parse(file);
        } catch (err) {
          throw new Error(err);
        }

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

export async function bundleFiles(configFile: string | undefined) {
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
async function upload(
  token: string,
  endpoint: string,
  formData: Record<string, any>
): Promise<PublishResponse> {
  const options = {
    method: 'PUT',
    url: endpoint,
    headers: {
      Authorization: `Bearer ${token}`,
      'Content-Type':
        typeof formData === 'object' ? 'application/json' : undefined,
    },
    formData,
  };

  return new Promise((resolve, reject) => {
    if (!namespace) {
      return reject(NAMESPACE_ERROR);
    }

    request(options, (error: any, response: any, body: any) => {
      if (error) {
        const errMessage =
          error.code === 'ETIMEDOUT' ? 'Publish request timed out' : error;
        return reject(errMessage);
      }

      if (fs.existsSync(TEMP_FILE)) {
        fs.unlinkSync(TEMP_FILE);
      }

      if (response.statusCode >= 400) {
        const message = body ? JSON.parse(body) : '';
        reject(
          new Error(
            `[${response.statusCode}] ${message.error}: ${message.results}`
          )
        );
      } else {
        const json = body ? JSON.parse(body) : {};
        resolve(json);
      }
    });
  });
}

async function update(
  token: string,
  endpoint: string,
  body: any
): Promise<PublishResponse> {
  const options = {
    method: 'PUT',
    url: endpoint,
    headers: {
      Authorization: `Bearer ${token}`,
    },
    json: true,
    body,
  };

  return new Promise((resolve, reject) => {
    if (!namespace) {
      return reject(NAMESPACE_ERROR);
    }

    request(options, (error: any, response: any, body: any) => {
      if (error) {
        const errMessage =
          error.code === 'ETIMEDOUT' ? 'Publish request timed out' : error;
        return reject(errMessage);
      }

      if (fs.existsSync(TEMP_FILE)) {
        fs.unlinkSync(TEMP_FILE);
      }

      if (response.statusCode >= 400) {
        const message = body ? JSON.parse(body) : '';
        reject(
          new Error(
            `[${response.statusCode}] ${message.error}: ${message.results}`
          )
        );
      } else {
        resolve(body);
      }
    });
  });
}

export async function publish(
  endpoint: string,
  options: PublishParams
): Promise<PublishResponse> {
  try {
    const { apiHost, authorizationEndpoint, dsApiHost, namespace } = config();
    const token = await authenticate(authorizationEndpoint);
    let path = endpoint;

    if (endpoint.includes(':')) {
      const compiler = compile(endpoint, {
        encode: encodeURIComponent,
      });
      path = compiler({ namespace });
    }
    const url = path.includes('/ds') ? dsApiHost + path : apiHost + path;

    const response = await update(token, url, options);
    return response;
  } catch (err) {
    throw new Error(err);
  }
}

export async function publishWithFile(
  endpoint: string,
  options: PublishParams
): Promise<PublishResponse> {
  try {
    const { apiHost, authorizationEndpoint, dsApiHost, namespace } = config();
    const token = await authenticate(authorizationEndpoint);
    let path = endpoint;

    if (endpoint.includes(':')) {
      const compiler = compile(endpoint, {
        encode: encodeURIComponent,
      });
      path = compiler({ namespace });
    }
    const url = path.includes('/ds') ? dsApiHost + path : apiHost + path;

    const response = await upload(token, url, options);
    return response;
  } catch (err) {
    throw new Error(err);
  }
}

export default publishWithFile;
