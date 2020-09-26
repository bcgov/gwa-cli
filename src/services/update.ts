import YAML from 'yaml';
import * as o2k from 'openapi-2-kong';

import { loadConfig } from './kong';
import { specState } from '../state/spec';
import { pluginsState } from '../state/plugins';
import { orgState } from '../state/org';

async function update(dir: string, files: string[]) {
  const [configFile, specFile] = files;

  if (configFile) {
    loadConfig(`${dir}/${configFile}`);
  }

  if (specFile) {
    const org = orgState.get();
    const result = await o2k.generate(
      `${dir}/${specFile}`,
      'kong-declarative-config',
      [org.name]
    );
    const [yamlDocs] = result.documents.map((d: any) => YAML.stringify(d));
  }
}

export default update;
