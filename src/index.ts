#!/usr/bin/env node
import dotenv from 'dotenv';
import chalk from 'chalk';
import { Command } from 'commander';
dotenv.config();

import run from './run';
import acl from './commands/acl';
import plugins from './commands/plugins';
import create from './commands/create';
import info from './commands/info';
import init from './commands/init';
import publish from './commands/publish';
import status from './commands/status';
import update from './commands/update';
// import validate from './commands/validate';
import { checkVersion } from './services/app';

const pkg = require('../package.json');

const program = new Command();
program.version(pkg.version);
// Refactored commands
program
  .addCommand(init)
  .addCommand(info)
  .addCommand(acl)
  .addCommand(publish)
  .addCommand(create)
  .addCommand(status);

// Commands to refactor

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

// program
//   .command('validate <input>')
//   .description('Validate a config file')
//   .action((input) => run(validate, input));

program
  .command('plugins [input]')
  .description('List all available plugins')
  .action((input) => run(plugins, input));

const main = async () => {
  try {
    const isValid = await checkVersion(pkg.version);
    if (isValid) {
      program.parse(process.argv);
    } else {
      console.log(
        chalk.bold
          .cyanBright`${chalk.yellow`[ Warning ]`} Your local version of APS CLI is out of date.`
      );
    }
  } catch (err) {
    throw err;
  }
};
try {
  main().catch(() => {
    process.exitCode = 1;
    console.log(
      chalk.bold.red`x Error`,
      'Unable to verify you have the latest version'
    );
  });
} catch (err) {
  process.exitCode = 1;
  console.log(chalk.bold.red`x Error`, err);
}
