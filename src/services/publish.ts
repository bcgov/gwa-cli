import { compile } from 'path-to-regexp';
import fs from 'fs';
import request from 'request';
import { extname, resolve } from 'path';

import authenticate from './auth';
import config from '../config';
import { getEndpoints } from './api';

const { apiHost, namespace } = config();
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
async function upload(
  token: string,
  endpoint: string,
  { configFile, dryRun }: PublishParams
): Promise<PublishResponse> {
  const value = await bundleFiles(configFile);
  const options = {
    method: 'PUT',
    url: endpoint,
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

    request(options, (error: Error, response: any, body: any) => {
      if (error) {
        reject(error);
      }
      const json = JSON.parse(body);

      if (response.statusCode >= 400) {
        reject(new Error(json.results || json.error));
      }

      if (fs.existsSync(TEMP_FILE)) {
        fs.unlinkSync(TEMP_FILE);
      }

      resolve(json);
    });
  });
}

// TODO combine this with api so there is only one entry point
export async function publish(
  endpoint: string,
  options: PublishParams
): Promise<PublishResponse> {
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
    const response = await upload(token, url, options);
    return response;
  } catch (err) {
    throw new Error(err);
  }
}

export default publish;
