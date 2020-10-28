import * as React from 'react';
import isString from 'lodash/isString';
import path from 'path';
import { render } from 'ink';
import validate from 'validate.js';

import ErrorView from '../views/error';
import { convertFile, convertRemote } from '../services/kong';
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
    if (input) {
      const isNotURL = validate.single(input, { url: true });
      const plugins = generatePluginTemplates(options.plugins, namespace);
      let output = null;

      if (isNotURL) {
        const file = path.resolve(cwd, input);
        const result = await importSpec(file);
        output = await convertRemote(result, namespace, plugins);
      } else {
        console.log('Fetching spec...');
        const result = await fetchSpec(input);
        output = await convertRemote(result, namespace, plugins);
      }

      if (isString(output)) {
        await exportConfig(output, options.outfile);
        console.log(`[DONE]: File ${options.outfile} generated`);
      }
    } else {
      const data = await loadPlugins(path.resolve(__dirname, '../../files'));
      initPluginsState(data);
      ui('/setup');
    }
  } catch (err) {
    render(<ErrorView text={err.message} title="New Config Failed" />);
    process.exit(0);
  }
}