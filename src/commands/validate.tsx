import get from 'lodash/get';
import path from 'path';
import validate from 'validate.js';

import constraints from '../validators';
import { validateConfig } from '../services/openapi';
import render from '../views/validate';

export default async function (input: string) {
  const cwd = process.cwd();

  try {
    const file = path.resolve(cwd, input);
    const validConfig = await validateConfig(file);

    const serviceErrors = validConfig.services.map((service: any) => {
      return service.plugins
        .map((plugin: any) => {
          const constraint: any | undefined = get(constraints, plugin.name);

          if (constraint) {
            const err = validate(plugin.config, constraint);

            if (err) {
              return {
                plugin: plugin.name,
                error: err,
              };
            }
          }

          return undefined;
        })
        .filter((err: any | undefined) => err !== undefined);
    });

    render(serviceErrors[0]);
  } catch (err) {
    console.error(err);
  }
}
