import has from 'lodash/has';
import head from 'lodash/head';
import minimist, { ParsedArgs } from 'minimist';
import path from 'path';

import { IAppContext } from './types';
import { loadConfig } from './services/kong';
import init from './views/init';
import render from './ui';
import help from './views/help';
import validate from './views/validate';

const p = require('../package.json');

const commands = ['edit', 'validate', 'init'];
const args: ParsedArgs = minimist(process.argv.slice(2));
const check = (args: ParsedArgs) => (...keys: string[]): boolean => {
  return keys.some((key) => has(args, key));
};

main(args);

function main(args: ParsedArgs) {
  const c = check(args);
  const version: string = p.version;
  const [command, ...rest] = args._;

  if (c('v', 'version')) {
    console.log(version);
    process.exit(1);
  }

  switch (command) {
    case 'edit':
      const file = rest[0];
      loadConfig(file);
      const config: IAppContext = {
        dir: args.dir,
        file,
        version,
      };
      render(config);
      break;
    case 'init':
      init(args);
      break;
    case 'validate':
      validate(rest[0]);
      break;
    default:
      help();
      process.exit(1);
      break;
  }
}
