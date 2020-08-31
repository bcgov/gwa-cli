import fetch from 'node-fetch';
import fs from 'fs';
import { resolve } from 'path';
import SwaggerParser from '@apidevtools/swagger-parser';
import YAML from 'yaml';
const o2k = require('openapi-2-kong');

import { specState } from '../state/spec';
import { pluginsState } from '../state/plugins';
import { orgState } from '../state/org';

export async function parseYaml(url: string, tag: string) {
  const res = await fetch(url);
  const json = await res.json();
  const api = await SwaggerParser.validate(json);
  let cache: any = [];
  const str = JSON.stringify(json, (key, value) => {
    if (typeof value === 'object' && value !== null) {
      if (cache.indexOf(value) !== -1) {
        // Circular reference found, discard key
        return;
      }
      // Store value in our collection
      cache.push(value);
    }
    return value;
  });
  cache = null;
  const result = await o2k.generateFromString(str, 'kong-declarative-config', [
    tag,
  ]);
  const [yamlDocs] = result.documents.map((d: any) => YAML.stringify(d));

  return YAML.parse(yamlDocs);
}

export function buildSpec(dir: string) {
  const spec = specState.get();
  const plugins = pluginsState.get();
  const org = orgState.get();
  const enabledPlugins = Object.values(plugins)
    .filter((p) => p.enabled)
    .map((p) => ({
      name: p.name,
      tags: [org.name],
      enabled: true,
      config: p.config,
    }));
  const configRef = JSON.parse(JSON.stringify(spec)); //TODO This is lazy, replace with a proper clone
  configRef.services[0].plugins = enabledPlugins;
  const specFile = YAML.stringify(configRef);
  fs.writeFileSync(resolve(dir + '/spec.yaml'), specFile);
}
