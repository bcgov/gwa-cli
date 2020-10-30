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
import { namespace } from '../config';
import ui from '../ui';

export default async function (input: string, options: any) {
  const cwd = process.cwd();

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
        output = await convertRemote(
          result,
          namespace,
          plugins,
          convertOptions
        );
      } else {
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
        await exportConfig(output, options.outfile);
        console.log(
          `${
            chalk.bold.green('✓ ') + chalk.bold('DONE ')
          } File ${chalk.italic.underline(options.outfile)} generated`
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
