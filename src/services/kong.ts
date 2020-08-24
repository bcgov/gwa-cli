const fs = require('fs');
const o2k = require('openapi-2-kong');
const { resolve } = require('path');
const SwaggerParser = require('@apidevtools/swagger-parser');
const YAML = require('yaml');

export async function parseYaml(json: any) {
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
    'geo',
  ]);
  const [yamlDocs] = result.documents.map((d: any) => YAML.stringify(d));

  return yamlDocs;
}
