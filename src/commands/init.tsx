import path from 'path';
import validate from 'validate.js';

import { exportConfig } from '../services/app';
import { fetchSpec, importSpec } from '../services/openapi';
import ui from '../ui';

export default async function (input: string, options: any) {
  const cwd = process.cwd();

  try {
    if (input) {
      const isNotURL = validate.single(input, { url: true });
      let output = null;

      if (isNotURL) {
        const file = path.resolve(cwd, input);
        output = await importSpec(file, options.team);
      } else {
        output = await fetchSpec(input, options.team);
      }
      await exportConfig(output, options.outfile);
      console.log(`[DONE]: File ${options.outfile} generated`);
    } else {
      ui();
    }
  } catch (err) {
    console.error(err);
  }
}
