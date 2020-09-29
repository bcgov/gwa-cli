import fs from 'fs';
import path from 'path';
import YAML from 'yaml';

//TODO: Add any importing or exporting of the config files in this module
export async function loadConfig(input: string): Promise<any> {
  const cwd = process.cwd();

  try {
    const file = await fs.promises.readFile(path.resolve(cwd, input), 'utf8');
    const json = YAML.parse(file);

    return json;
  } catch (err) {
    console.log(err);
  }
}

export function parseConfig(json: any) {
  const service = json.services[0];
  const name = service.name;
  const team = service.tags.slice(-1)[0];
  const host = service.host || service.url;
  const plugins = service.plugins;

  return {
    name,
    team,
    host,
    plugins,
  };
}

export async function exportConfig(
  output: string,
  outfile: string
): Promise<any> {
  const cwd = process.cwd();

  try {
    const specFile = YAML.stringify(output);
    await fs.promises.writeFile(
      path.resolve(cwd, outfile),
      specFile.replace('|', '')
    );
  } catch (err) {
    console.error(err);
  }
}

