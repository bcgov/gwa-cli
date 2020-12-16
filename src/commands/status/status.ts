import { Command } from 'commander';

import render from './renderer';

export const actionHandler = () => {
  render();
};

const status = new Command('status');

status.description('Check the status of your configs').action(actionHandler);

export default status;
