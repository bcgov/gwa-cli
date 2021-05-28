import chalk from 'chalk';
import { Command } from 'commander';

import { checkForEnvFile } from '../../services/app';
import render from './renderer';

export const actionHandler = () => {
  if (checkForEnvFile()) {
    render();
  } else {
    console.log(
      chalk.cyanBright`${chalk.bold
        .yellow`[ Warning ]`} No .env file found. Run ${chalk.white`gwa init`} to start a new configuration environment.`
    );
  }
};

const info = new Command('info');

info.description('Print info about your env').action(actionHandler);

export default info;
