import { program } from 'commander';

import run from './run';
import plugins from './commands/plugins';
import init from './commands/init';

const pkg = require('../package.json');

const main = async () => {
  program
    .command('new [input]')
    .option('-t, --team <team>', 'The team you wish to register')
    .option(
      '-o, --outfile <output>',
      'An OpenAPI spec JSON file on your local machine'
    )
    .description('Initialize a config file in the current directory')
    .action((input, options) => run(init, input, options));

  program
    .command('edit <input>')
    .description('Edit a config file')
    .action(() => run('edit'));

  program
    .command('update <input>')
    .description('Update a config with new OpenAPI specs')
    .option('-u, --url <url>', 'The URL of a OpenAPI spec JSON file')
    .option(
      '-f, --file <file>',
      'An OpenAPI spec JSON file on your local machine'
    )
    .action(() => run('update'));

  program
    .command('validate <input>')
    .description('Validate a config file')
    .action(() => run('validate'));

  program
    .command('plugins [input]')
    .description('List all available plugins')
    .action((...args) => run(plugins, ...args));

  program.version(pkg.version, '-v, --version');
  await program.parseAsync(process.argv);
};

main();
