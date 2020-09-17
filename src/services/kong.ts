import fetch from 'node-fetch';
import fs from 'fs';
import merge from 'deepmerge';
import mapValues from 'lodash/mapValues';
import { resolve } from 'path';
import SwaggerParser from '@apidevtools/swagger-parser';
import YAML from 'yaml';
import * as o2k from 'openapi-2-kong';

import { specState } from '../state/spec';
import { pluginsState, encryptedValues } from '../state/plugins';
import { orgState } from '../state/org';

export async function parseYaml(url: string, tag: string) {
  const res = await fetch(url);
  const json = await res.json();
  await SwaggerParser.validate(json);
  let cache: any = [];
  const str = JSON.stringify(json, (_, value) => {
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

export function loadConfig(file: string) {
  fs.readFile(file, 'utf8', (err, data) => {
    if (err) {
      return;
    }

    const result = YAML.parse(data);
    const name = result.services[0].tags.slice(-1)[0];
    const host = result.services[0].host;
    const plugins = result.services[0].plugins;
    specState.set(result);
    orgState.set((prev) => ({
      ...prev,
      name,
      host,
      file,
    }));
    plugins.forEach((plugin: any) => {
      if (pluginsState.get()[plugin.name]) {
        pluginsState.set((prev) => {
          return {
            ...prev,
            [plugin.name]: {
              ...prev[plugin.name],
              data: merge(prev[plugin.name].data, plugin),
              enabled: plugin.enabled,
            },
          };
        });
      }
    });
  });
}

export async function encryptValue(value: string) {
  try {
    const req = await fetch(
      `https://gwa-qwzrwc-dev.pathfinder.gov.bc.ca/encrypt?v=${value}`
    );
    const text = await req.text();
    return text;
  } catch (err) {
    throw err;
  }
}

export function buildSpec(
  dir: string,
  file: string | null = 'spec.yaml'
): void {
  const spec = specState.get();
  const plugins = pluginsState.get();
  const org = orgState.get();
  const enabledPlugins = Object.values(plugins)
    .filter((p) => p.data.enabled)
    .map((p) => ({
      name: p.data.name,
      tags: [org.name],
      enabled: true,
      config: mapValues(p.data.config, (v, k) => v),
    }));
  const configRef = JSON.parse(JSON.stringify(spec)); //TODO This is lazy, replace with a proper clone
  configRef.services[0].plugins = enabledPlugins;
  const specFile = YAML.stringify(configRef);
  fs.writeFileSync(resolve(`${dir}/${file}`), specFile);
}
