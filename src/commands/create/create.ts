import chalk from 'chalk';
import { Command } from 'commander';

import { makeConfigFile } from './create-actions';
import render from './renderer';

const actionHandler = async (input: string | undefined, options: any) => {
  if (input) {
    try {
      console.log(chalk.italic('Fetching spec...'));
      const result = await makeConfigFile(input, options);

      console.log(
        `${
          chalk.bold.green('✓ ') + chalk.bold('DONE ')
        } File ${chalk.italic.underline(result)} generated`
      );
    } catch (err) {
      process.exitCode = 1;
      console.log(chalk.bold.red`x Error:` + ` ${err.message}`);
    }
  } else {
    render(makeConfigFile);
  }
};

const create = new Command('new');
create
  .arguments('[input]')
  //.option('--service <service>', "The service's name")
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
