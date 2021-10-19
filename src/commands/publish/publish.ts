import { Command } from 'commander';

import render from './renderer';

type PublishOptions = {
  body: string;
  content?: string;
  verbose?: boolean;
};

export const actionHandler = (action: string, options: PublishOptions) => {
  if (/(content|dataset|product|issuer)/.test(action)) {
    render(action, options);
  }
};

const publish = new Command('publish');

publish
  .description(
    'Publish to DS API. Available commands are content, datasets, issuers and products'
  )
  .arguments('[action]')
  .option('-b,--body <body>', 'YAML file to convert to JSON')
  .option('-c,--content [content]', 'Content to add to body')
  .option('--verbose')
  .option('--debug')
  .action(actionHandler);

export default publish;
