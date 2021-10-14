import chalk from 'chalk';
import fetch from 'node-fetch';
import flatten from 'lodash/flatten';
import fs from 'fs';
import path from 'path';
import validate from 'validate.js';
import YAML from 'yaml';
import type { InitOptions } from '../types';

export async function checkVersion(
  pkgVersion: string
): Promise<boolean | string> {
  try {
    const res = await fetch(
      'https://api.github.com/repos/bcgov/gwa-cli/releases/latest'
    );
    const json = await res.json();
    const currentVersion = json.tag_name.replace('v', '');

    if (pkgVersion < currentVersion) {
      return currentVersion;
    }

    return true;
  } catch (err) {
    throw err;
  }
}

export function checkForEnvFile() {
  return fs.existsSync('.env');
}

export function checkForApiVersion() {
  if (checkForEnvFile() && !process.env.API_VERSION) {
    console.log('');
    console.log(
      chalk.cyanBright`${chalk.bold
        .yellow`[ Warning ]`} Your ${chalk.underline`API_VERSION`} has not be set.`
    );
    console.log(
      'All requests will default to v2. ACL commands are not available.'
    );
    console.log('');
  }
}

export async function makeEnvFile(options: InitOptions): Promise<string> {
  const rules = {
    namespace: {
      presence: { allowEmpty: false },
      length: { minimum: 5, maximum: 15 },
      format: {
        pattern: '^[a-z][a-z0-9-]{4,14}$',
        flags: 'gi',
        message: 'can only contain a-z, 0-9 and dashes',
      },
    },
    dataCenter: {
      presence: { allowEmpty: true },
    },
    clientId: {
      presence: { allowEmpty: false },
    },
    clientSecret: {
      presence: { allowEmpty: false },
    },
    apiVersion: {
      format: {
        pattern: '[1-2]+',
        message: 'only versions 1 and 2 are available',
      },
    },
  };

  try {
    const errors = validate(options, rules);

    if (errors) {
      const errorMessage = flatten(Object.values(errors)).join('\n');

      throw errorMessage;
    }

    let data = `GWA_NAMESPACE=${options.namespace}
CLIENT_ID=${options.clientId}
CLIENT_SECRET=${options.clientSecret}
GWA_ENV=${options.env}
API_VERSION=${options.apiVersion ?? '2'}
`;
    if (options.dataCenter) {
      data += `DATA_CENTER=${options.dataCenter}`;
    }
    await fs.promises.writeFile('.env', data);
    return '.env file successfully generated';
  } catch (err) {
    throw new Error(`Unable to write file ${err}`);
  }
}

export async function loadConfig(input: string): Promise<any> {
  const cwd = process.cwd();

  try {
    const file = await fs.promises.readFile(path.resolve(cwd, input), 'utf8');
    const json = YAML.parse(file);

    return json;
  } catch (err) {
    throw new Error(err);
  }
}

export async function saveConfig(
  output: string,
  outfile: string
): Promise<any> {
  const cwd = process.cwd();

  try {
    await fs.promises.writeFile(path.resolve(cwd, outfile), output);
  } catch (err) {
    throw new Error(err);
  }
}
