import * as o2k from '../lib/openapi-2-kong';
import YAML from 'yaml';

function extract(result: any): string {
  const [doc] = result.documents.map((d: any) => YAML.stringify(d));
  return doc;
}

export async function convertFile(
  file: string,
  namespace: string,
  plugins?: string[]
): Promise<string> {
  try {
    const result = await o2k.generate(file, 'kong-declarative-config', [
      `ns.${namespace}`,
    ]);
    const document = extract(result);

    if (plugins && plugins.length > 0) {
      return addPluginsToSpec(document, plugins);
    }

    return document;
  } catch (err) {
    throw err.message;
  }
}

export type GenerateConfigOptions = {
  input: string;
  namespace: string;
  plugins: string[];
  options: {
    routeHost: string;
    serviceUrl: string;
  };
};

export async function generateConfig({
  input,
  namespace,
  plugins,
  options,
}: GenerateConfigOptions): Promise<string> {
  try {
    const result = await o2k.generateFromString(
      input,
      'kong-declarative-config',
      [`ns.${namespace}`],
      options
    );
    const document = extract(result);

    if (plugins && plugins.length > 0) {
      return addPluginsToSpec(document, plugins);
    }

    return document;
  } catch (err) {
    throw err.message;
  }
}

export function addPluginsToSpec(config: string, plugins: any[]): string {
  const json = YAML.parse(config);
  const services = json.services.map((service: any) => ({
    ...service,
    plugins,
  }));
  return YAML.stringify({ ...json, services });
}
