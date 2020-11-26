import compact from 'lodash/compact';
import isString from 'lodash/isString';
import path from 'path';
import validate from 'validate.js';

import { convertRemote, ImportOptions } from '../../services/kong';
import { exportConfig } from '../../services/app';
import { fetchSpec, importSpec } from '../../services/openapi';
import { initPluginsState, generatePluginTemplates } from '../../state/plugins';
import config from '../../config';

export const parsePlugins = (plugins: string | string[], namespace: string) => {
  const requestedPlugins = isString(plugins)
    ? compact(plugins.split(/,|\s/g))
    : plugins;
  return generatePluginTemplates(requestedPlugins, namespace);
};

export const makeConfigFile = async (
  input: string,
  options: {
    outfile: string;
    plugins: string[];
    routeHost: string;
    serviceUrl: string;
  }
): Promise<string> => {
  const { namespace } = config();
  const cwd = process.cwd();
  let outfile = options.outfile;

  try {
    await initPluginsState();

    const { routeHost, serviceUrl } = options;
    const isNotURL = validate.single(input, { url: true });
    const plugins = parsePlugins(options.plugins, namespace);
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
