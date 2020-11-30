import compact from 'lodash/compact';
import flow from 'lodash/flow';
import isString from 'lodash/isString';
import path from 'path';

import { generateConfig, GenerateConfigOptions } from '../../services/kong';
import { exportConfig } from '../../services/app';
import { fetchSpec, importSpec } from '../../services/openapi';
import { initPluginsState, generatePluginTemplates } from '../../state/plugins';
import { isLocalInput, makeOutputFilename } from '../../services/utils';
import config from '../../config';

export const importInput = async (
  options: GenerateConfigOptions
): Promise<GenerateConfigOptions> => {
  let input;
  const isLocal = isLocalInput(options.input);

  try {
    if (isLocal) {
      const file = path.resolve(process.cwd(), options.input);
      input = await importSpec(file);
    } else {
      input = await fetchSpec(options.input);
    }

    return {
      ...options,
      input,
    };
  } catch (err) {
    throw new Error(err);
  }
};

export const parseOptions = (
  options: GenerateConfigOptions
): GenerateConfigOptions => {
  const requestedPlugins = isString(options.plugins)
    ? compact(options.plugins.split(/,|\s/g))
    : options.plugins;
  const plugins = generatePluginTemplates(requestedPlugins, options.namespace);

  return {
    ...options,
    plugins,
  };
};

export const processInput = flow([parseOptions, importInput]);

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
  const { routeHost, serviceUrl } = options;
  const outfile = makeOutputFilename(input, options.outfile);

  try {
    await initPluginsState();
    const result = await processInput({
      input,
      namespace,
      plugins: options.plugins,
      options: {
        routeHost,
        serviceUrl,
      },
    });
    const output = await generateConfig(result);

    await exportConfig(output, outfile);
    return outfile;
  } catch (err) {
    throw new Error(err);
  }
};
