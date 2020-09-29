import fs from 'fs';
import path from 'path';
import YAML from 'yaml';

//TODO: Add any importing or exporting of the config files in this module
export async function exportConfig(output: string, outfile: string) {
  const cwd = process.cwd();

  try {
    const specFile = YAML.stringify(output);
    await fs.promises.writeFile(path.resolve(cwd, outfile), specFile);
  } catch (err) {
    console.error(err);
  }
}

