import fs from 'fs';
import { basename, extname, resolve } from 'path';
import YAML from 'yaml';

import { PluginsResult } from '../types';

/**
 * Crawls a directory for YAML config files and outputs a dictionary of results
 *
 * Entries are comprised of the meta data of the plugin (basic details) and the config is the possible values
 */
export async function loadPlugins(path: string): Promise<PluginsResult> {
  const result: PluginsResult = {};

  try {
    const dir = await fs.promises.opendir(path);

    for await (const dirent of dir) {
      const extension = extname(dirent.name);

      if (/\.(yaml|yml)/.test(extension)) {
        const id = basename(dirent.name, extension);
        const file = await fs.promises.readFile(
          resolve(path, dirent.name),
          'utf8'
        );
        // Files are divided into 2 documents, the meta and the config fields
        const [meta, config] = file.split('---').map((d) => YAML.parse(d));
        result[id] = {
          meta: {
            id,
            ...meta,
          },
          config,
        };
      }
    }
  } catch (err) {
    console.error(err);
  }

  return result;
}

export function convert(config: any): string {}
