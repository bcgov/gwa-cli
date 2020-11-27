import fs from 'fs';
import isString from 'lodash/isString';
import path from 'path';
import validate from 'validate.js';

import { convertRemote, ImportOptions } from '../../services/kong';
import { exportConfig } from '../../services/app';
import { fetchSpec, importSpec } from '../../services/openapi';
import { loadPlugins } from '../../services/plugins';
import { initPluginsState, generatePluginTemplates } from '../../state/plugins';
import config from '../../config';

const { namespace } = config();

export const makeConfigFile = async (
  input: string,
  options: {
    outfile: string;
    plugins: string[];
    routeHost: string;
    serviceUrl: string;
  }
): Promise<string> => {
  const cwd = process.cwd();
  let outfile = options.outfile;

  try {
    const data = await loadPlugins();
    initPluginsState(data);

    const { routeHost, serviceUrl } = options;
    const isNotURL = validate.single(input, { url: true });
    const requestedPlugins = isString(options.plugins)
      ? options.plugins.split(/,|\s/g)
      : options.plugins;
    const plugins = generatePluginTemplates(requestedPlugins, namespace);
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

      output = await convertRemote(result, namespace, plugins, convertOptions);
    } else {
      if (!outfile) {
        throw new Error('An --outfile must be set');
      }

      const result = await fetchSpec(input);
      output = await convertRemote(result, namespace, plugins, convertOptions);
    }

    if (!isString(output)) {
      throw new Error('No output generated');
    }

    await exportConfig(output, outfile);
    return outfile;
  } catch (err) {
    throw new Error(err);
  }
};
