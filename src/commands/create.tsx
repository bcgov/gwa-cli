import * as React from 'react';
import chalk from 'chalk';
import isString from 'lodash/isString';
import path from 'path';
import { render } from 'ink';
import validate from 'validate.js';

import ErrorView from '../views/error';
import { convertFile, convertRemote, ImportOptions } from '../services/kong';
import { exportConfig } from '../services/app';
import { fetchSpec, importSpec } from '../services/openapi';
import { loadPlugins } from '../services/plugins';
import { generatePluginTemplates } from '../state/plugins';
import { initPluginsState } from '../state/plugins';
import config from '../config';
import ui from '../ui';

const {
  clientId,
  clientSecret,
  apiHost,
  authorizationEndpoint,
  namespace,
} = config();

export default async function (input: string, options: any) {
  const cwd = process.cwd();
  let outfile = options.outfile;

  try {
    const data = await loadPlugins(path.resolve(__dirname, '../../files'));
    initPluginsState(data);

    if (input) {
      const { routeHost, serviceUrl } = options;
      const isNotURL = validate.single(input, { url: true });
      const plugins = generatePluginTemplates(options.plugins, namespace);
      const convertOptions: ImportOptions = {
        routeHost,
        serviceUrl,
      };
      let output = null;

      if (isNotURL) {
        const file = path.resolve(cwd, input);
        const result = await importSpec(file);

        if (!options.outfile) {
          outfile = input.replace(/json$/i, 'yaml');
        }

        output = await convertRemote(
          result,
          namespace,
          plugins,
          convertOptions
        );
      } else {
        if (!outfile) {
          throw new Error('An --outfile must be set');
        }

        console.log(chalk.italic('Fetching spec...'));
        const result = await fetchSpec(input);
        output = await convertRemote(
          result,
          namespace,
          plugins,
          convertOptions
        );
      }

      if (isString(output)) {
        await exportConfig(output, outfile);
        console.log(
          `${
            chalk.bold.green('✓ ') + chalk.bold('DONE ')
          } File ${chalk.italic.underline(outfile)} generated`
        );
      }
    } else {
      ui('/setup');
    }
  } catch (err) {
    process.exitCode = 1;
    render(<ErrorView text={err.message} title="New Config Failed" />);
  }
}
