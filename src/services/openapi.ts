import fetch from 'node-fetch';
import fs from 'fs';
import SwaggerParser from '@apidevtools/swagger-parser';
import path from 'path';
import YAML from 'yaml';

export async function validateConfig(file: string): Promise<any> {
  try {
    const contents = await fs.promises.readFile(file, 'utf8');
    const json = YAML.parse(contents);
    // TODO: find a way to validate routes after converted to kong
    // await SwaggerParser.validate(json);
    return json;
  } catch (err) {
    throw new Error(err);
  }
}

// Read and validate a local OpenAPI JSON file
export async function importSpec(file: string) {
  try {
    const fileType = path.extname(file);
    const contents = await fs.promises.readFile(file, 'utf8');
    const json =
      fileType === '.yaml' ? YAML.parse(contents) : JSON.parse(contents);
    await SwaggerParser.validate(json);
    const result = parser(json);
    return result;
  } catch (err) {
    throw new Error(err);
  }
}

// Fetch and validate a hosted OpenAPI spec
export async function fetchSpec(url: string) {
  try {
    const res = await fetch(url);
    let json;
    if (/\.yaml$/.test(url)) {
      const yamlResponse = await res.text();
      json = YAML.parse(yamlResponse);
    } else {
      json = await res.json();
    }
    await SwaggerParser.validate(json);
    const result = parser(json);
    return result;
  } catch (err) {
    throw new Error(err);
  }
}

function parser(json: any) {
  let cache: any = [];
  // Need to remove any circular references
  const result = JSON.stringify(json, (_, value) => {
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
  return result;
}
