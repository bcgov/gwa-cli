import { Command } from 'commander';

import render from './renderer';

type PublishOptions = {
  dryRun: boolean;
};

export const actionHandler = (input: string, options: PublishOptions) => {
  render({
    configFile: input,
    dryRun: Boolean(options.dryRun).toString(),
  });
};

const publish = new Command('publish-gateway');

publish
  .alias('pg')
  .description('Publish gateway config')
  .arguments('[input]')
  .option('--dry-run', 'Enable dry run')
  .option('--debug')
  .action(actionHandler);

export default publish;
