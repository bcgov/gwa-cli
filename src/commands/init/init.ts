import chalk from 'chalk';
import { Command } from 'commander';
import isEmpty from 'lodash/isEmpty';
import pick from 'lodash/pick';

import { checkForEnvFile, makeEnvFile } from '../../services/app';
import prompts from './prompts';
import { render } from '../../views/create-env';
import type { InitOptions } from '../../types';

export const actionHandler = (options: InitOptions) => {
  const initArgs = pick(options, ['namespace', 'clientId', 'clientSecret']);
  const envArgs = pick(options, ['dev', 'test', 'prod']);
  const env = Object.keys(envArgs)[0] ?? 'test';

  if (checkForEnvFile()) {
    throw 'You already have an .env file in this project';
  }

  // Empty args indicates we launch the prompt form instead
  if (isEmpty(initArgs)) {
    render(prompts, env);
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
        process.exitCode = 1;

        if (options.debug) {
          console.error(err);
        }
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
  .option('--client-id <clientId>', 'The Service Account Client ID')
  .option('--client-secret <clientSecret>', 'The Service Account Client Secret')
  .option('--debug', 'Show stack traces on error. Useful for debugging.')
  .action(actionHandler);

export default init;
