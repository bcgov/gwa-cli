import path from 'path';
import { program } from 'commander';

//import { IAppContext } from './types';
//import { loadConfig } from './services/kong';
//import init from './views/init';
//import render from './ui';
//import update from './services/update';
//import validate from './views/validate';
import plugins from './domains/plugins';

const pkg = require('../package.json');

const main = () => {
  program
    .command('new')
    .description('Initialize a config file in the current directory');

  program.command('edit <input>').description('Edit a config file');

  program
    .command('update <input>')
    .description('Update a config with new OpenAPI specs')
    .option('-u, --url <url>', 'The URL of a OpenAPI spec JSON file')
    .option(
      '-f, --file <file>',
      'An OpenAPI spec JSON file on your local machine'
    )
    .action((input, options) => {
      console.log(input, options.url, options.file);
    });

  program
    .command('validate <input>')
    .description('Validate a config file')
    .action((input) => console.log('input', input));

  program
    .command('list')
    .alias('ls')
    .description('List all available plugins')
    .action((cmd, options) => plugins(path.resolve(__dirname, '../files')));

  program.version(pkg.version, '-v, --version');
  program.parse(process.argv);
};

main();
