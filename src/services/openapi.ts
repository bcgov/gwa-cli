import fetch from 'node-fetch';
import fs from 'fs';
import * as o2k from 'openapi-2-kong';
import SwaggerParser from '@apidevtools/swagger-parser';
import YAML from 'yaml';

function extract(result: any) {
  const [doc] = result.documents.map((d: any) => YAML.stringify(d));
  return doc;
}

export async function importSpec(file: string, tag: string) {
  try {
    const contents = await fs.promises.readFile(file, 'utf8');
    const json = JSON.parse(contents);
    await SwaggerParser.validate(json);
    const result = await o2k.generate(file, 'kong-declarative-config', [tag]);
    return extract(result);
  } catch (err) {
    console.error(err);
  }
}

export async function fetchSpec(url: string, tag: string) {
  try {
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
    const result = await o2k.generateFromString(
      str,
      'kong-declarative-config',
      [tag]
    );
    return extract(result);
  } catch (err) {
    console.error(err);
  }
}
