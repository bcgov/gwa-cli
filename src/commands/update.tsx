import isString from 'lodash/isString';
import path from 'path';

import { exportConfig, loadConfig } from '../services/app';
import { fetchSpec, importSpec } from '../services/openapi';
import { generateConfig } from '../services/kong';

export default async function (input: string, options: any) {
  const cwd = process.cwd();

  try {
    const draft = await loadConfig(input);
    const plugins = draft.services.map((service: any) => service.plugins);
    let output = null;

    if (options.url) {
      console.log('Fetching spec...');
      const result = await fetchSpec(options.url);
      output = await generateConfig({
        input: result,
        namespace: options.team,
        plugins,
        options: {
          routeHost: options.routeHost,
          serviceUrl: options.serviceUrl,
        },
      });
    } else if (options.file) {
      const file = path.resolve(cwd, options.file);
      const result = await importSpec(file);
      output = await generateConfig({
        input: result,
        namespace: options.team,
        plugins,
        options: {
          routeHost: options.routeHost,
          serviceUrl: options.serviceUrl,
        },
      });
    }

    if (isString(output)) {
      await exportConfig(output, input);
      console.log(`[DONE]: File ${input} updated`);
    }
  } catch (err) {
    console.error(err);
  }
}
