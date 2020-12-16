import { basename, extname } from 'path';
import fs from 'fs';
import path from 'path';
import YAML from 'yaml';

import { PluginsResult } from '../types';

/**
 * Crawls a directory for YAML config files and outputs a dictionary of results
 *
 * Entries are comprised of the meta data of the plugin (basic details) and the config is the possible values
 */
export async function loadPlugins(): Promise<PluginsResult> {
  const result: PluginsResult = {};
  const folder = path.join(__dirname, '../../files');
  try {
    const files = fs.readdirSync(folder);

    // files object contains all files names
    // log them on console
    files.forEach((file) => {
      const extension = extname(file);

      if (/\.(yaml|yml)/.test(extension)) {
        const id = basename(file, extension);
        const configFile = fs.readFileSync(
          path.join(__dirname, `../../files/${file}`),
          'utf8'
        );
        // Files are divided into 2 documents, the meta and the config fields
        const [meta, config] = configFile
          .split('---')
          .map((d) => YAML.parse(d));
        result[id] = {
          meta: {
            id,
            ...meta,
          },
          config,
        };
      }
    });
  } catch (err) {
    throw new Error(err);
  }

  return result;
}
