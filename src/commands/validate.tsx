import get from 'lodash/get';
import path from 'path';

import schemas from '../validators';
import { validateConfig } from '../services/openapi';
import render from '../views/validate';

export default async function (input: string) {
  const cwd = process.cwd();

  try {
    const file = path.resolve(cwd, input);
    const validConfig = await validateConfig(file);

    const serviceErrors: any[] = validConfig.services.map((service: any) => {
      return service.plugins
        .map((plugin: any) => {
          const schema: any | undefined = get(schemas, plugin.name);

          if (schema) {
            const { error, value } = schema.validate(plugin.config);

            if (error) {
              return {
                plugin: plugin.name,
                error,
                value,
              };
            }
          }

          return undefined;
        })
        .filter((err: any | undefined) => err !== undefined);
    });
    render(serviceErrors[0]);
  } catch (err) {
    throw new Error(err);
  }
}
