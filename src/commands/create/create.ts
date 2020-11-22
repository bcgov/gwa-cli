import chalk from 'chalk';
import { Command } from 'commander';

import { makeConfigFile } from './create-actions';
import render from './renderer';

export const actionHandler = (input: string | undefined, options: any) => {
  if (input) {
    console.log(chalk.italic('Fetching spec...'));
    makeConfigFile(input, options).then(
      (result) => {
        console.log(
          `${
            chalk.bold.green('✓ ') + chalk.bold('DONE ')
          } File ${chalk.italic.underline(result)} generated`
        );
      },
      (err) => {
        process.exitCode = 1;
        console.error(chalk.bold.red`x Error:` + ` ${err.message}`);
      }
    );
  } else {
    render();
  }
};

const create = new Command('new');
create
  .arguments('[input]')
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
  .action(actionHandler);

export default create;
