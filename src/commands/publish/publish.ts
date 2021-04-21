import { Command } from 'commander';

import render from './renderer';

type PublishOptions = {
  dryRun?: boolean;
  verbose?: boolean;
};

export const actionHandler = (input: string, options: PublishOptions = {}) => {
  render({
    configFile: input,
    dryRun: Boolean(options.dryRun).toString(),
    verbose: options.verbose,
  });
};

const publish = new Command('publish-gateway');

publish
  .alias('pg')
  .description('Publish gateway config')
  .arguments('[input]')
  .option('--dry-run', 'Enable dry run')
  .option('--verbose')
  .option('--debug')
  .action(actionHandler);

export default publish;
