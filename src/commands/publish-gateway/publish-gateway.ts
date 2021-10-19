import { Command } from 'commander';

import render from './renderer';

type PublishGatewayOptions = {
  dryRun?: boolean;
  verbose?: boolean;
};

export const actionHandler = (
  input: string,
  options: PublishGatewayOptions = {}
) => {
  render({
    configFile: input,
    dryRun: Boolean(options.dryRun).toString(),
    verbose: options.verbose,
  });
};

const publishGateway = new Command('publish-gateway');

publishGateway
  .alias('pg')
  .description('Publish gateway config')
  .arguments('[input]')
  .option('--dry-run', 'Enable dry run')
  .option('--verbose')
  .option('--debug')
  .action(actionHandler);

export default publishGateway;
