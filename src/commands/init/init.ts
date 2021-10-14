import chalk from 'chalk';
import { Command } from 'commander';
import isEmpty from 'lodash/isEmpty';
import pick from 'lodash/pick';

import { checkForEnvFile, makeEnvFile } from '../../services/app';
import render from './renderer';
import type { InitOptions } from '../../types';

export const actionHandler = (options: InitOptions) => {
  const initArgs = pick(options, [
    'namespace',
    'clientId',
    'clientSecret',
    'apiVersion',
  ]);
  const envArgs = pick(options, ['dev', 'test', 'prod']);
  const env = Object.keys(envArgs)[0] ?? 'test';

  if (checkForEnvFile()) {
    throw 'You already have an .env file in this project';
  }

  // Empty args indicates we launch the prompt form instead
  if (isEmpty(initArgs)) {
    render(env);
  } else {
    makeEnvFile({
      ...options,
      env,
    }).then(
      (message: string) => {
        console.log(chalk.green.bold('Success'), message);
      },
      (err) => {
        console.error(chalk.red.bold('x Error'), 'Unable to create .env file');
        if (err) {
          console.error('Reason:');
          console.error(err.message);
        }
        process.exitCode = 1;

        if (options.debug) {
          console.error(err);
        }

        process.exit();
      }
    );
  }
};

const init = new Command('init');

init
  .option(
    '--namespace <namespace>',
    'Represents a collections of Kong Services and Routes'
  )
  .option('-D, --dev', 'Dev environment')
  .option('-P, --prod', 'Set environment to production')
  .option('-T, --test', 'Set environment to test (default)')
  .option('--data-center [dc]', 'Target a particular data center', '')
  .option('--client-id <clientId>', 'The Service Account Client ID')
  .option('--client-secret <clientSecret>', 'The Service Account Client Secret')
  .option('--debug', 'Show stack traces on error. Useful for debugging.')
  .option(
    '--api-version <apiVersion>',
    'Show stack traces on error. Useful for debugging.'
  )
  .action(actionHandler);

export default init;
