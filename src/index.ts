import has from 'lodash/has';
import head from 'lodash/head';
import minimist, { ParsedArgs } from 'minimist';
import path from 'path';

import { IAppContext } from './types';
import { loadConfig } from './services/kong';
import render from './ui';

const p = require('../package.json');

const args: ParsedArgs = minimist(process.argv.slice(2));
const check = (args: ParsedArgs) => (...keys: string[]): boolean => {
  return keys.some((key) => has(args, key));
};

main(args);

function main(args: ParsedArgs) {
  const c = check(args);
  const version: string = p.version;
  const file = head(args._) || '';

  if (file) {
    // TODO: Make this async
    loadConfig(file);
  }

  if (c('v', 'version')) {
    console.log();
    process.exit(1);
  } else if (c('c', 'check')) {
    console.log('validate');
    process.exit(1);
  } else if (args._.includes('init')) {
    if (c('org')) {
      console.log(`Generating ${args.org}...`);
      setTimeout(() => {
        console.clear();
        console.log('Done');
      }, 3000);
    }
  } else {
    const config: IAppContext = {
      dir: args.dir,
      file,
      version,
    };
    render(config);
  }
}
