#!/usr/bin/env node
import dotenv from 'dotenv';
import chalk from 'chalk';
import { Command } from 'commander';
dotenv.config();

import run from './run';
import acl from './commands/acl';
import plugins from './commands/plugins';
import create from './commands/create';
import edit from './commands/edit';
import init from './commands/init';
import publish from './commands/publish';
import update from './commands/update';
import validate from './commands/validate';

const pkg = require('../package.json');

const program = new Command();
program.version(pkg.version);
// Refactored commands
program.addCommand(init).addCommand(acl).addCommand(publish);

// Commands to refactor
program
  .command('new [input]')
  .option('--service <service>', "The service's name")
  .option('--route-host <routeHost>', "Generally a server's URL")
  .option(
    '--service-url <serviceUrl>',
    'The URL the server will be accessable through'
  )
  .option(
    '-p, --plugins [plugins...]',
    'Any starter plugins you would like to include'
  )
  .option(
    '-o, --outfile <output>',
    'An OpenAPI spec JSON file on your local machine'
  )
  .description(
    'Initialize a config file in the current directory. The input file must be an OpenAPI JSON file or URL'
  )
  .option('--debug')
  .action((input, options) => run(create, input, options));

program
  .command('edit <input>')
  .description('Edit a config file')
  .action((input) => run(edit, input));

program
  .command('update <input>')
  .description('Update a config with new OpenAPI specs')
  .option('-u, --url [url]', 'The URL of a OpenAPI spec JSON file')
  .option(
    '-f, --file [file]',
    'An OpenAPI spec JSON file on your local machine'
  )
  .option('--debug')
  .action((input, options) => run(update, input, options));

program
  .command('validate <input>')
  .description('Validate a config file')
  .action((input) => run(validate, input));

program
  .command('plugins [input]')
  .description('List all available plugins')
  .action((input) => run(plugins, input));

try {
  program.parse(process.argv);
} catch (err) {
  process.exitCode = 1;
  console.log(chalk.bold.red`x Error`, err);
}
