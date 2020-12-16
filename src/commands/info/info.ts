import { Command } from 'commander';

import render from './renderer';

export const actionHandler = () => {
  render();
};

const info = new Command('info');

info.description('Print info about your env').action(actionHandler);

export default info;
