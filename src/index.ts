import { program } from 'commander';

import run from './run';
import plugins from './commands/plugins';
import init from './commands/init';
import edit from './commands/edit';
import publish from './commands/publish';
import update from './commands/update';
import validate from './commands/validate';

const pkg = require('../package.json');

const main = async () => {
  program
    .command('new [input]')
    .option('-t, --team <team>', 'The team you wish to register')
    .option(
      '-p, --plugins [plugins...]',
      'Any starter plugins you would like to include'
    )
    .option(
      '-o, --outfile <output>',
      'An OpenAPI spec JSON file on your local machine'
    )
    .description(
      'Initialize a config file in the current directory. The input file must be an OpenAPI JSON file'
    )
    .action((input, options) => run(init, input, options));

  program
    .command('edit <input>')
    .description('Edit a config file')
    .action((input) => run(edit, input));

  program
    .command('update <input>')
    .description('Update a config with new OpenAPI specs')
    .option('-t, --team <team>', 'You must declare the team for this spec')
    .option('-u, --url [url]', 'The URL of a OpenAPI spec JSON file')
    .option(
      '-f, --file [file]',
      'An OpenAPI spec JSON file on your local machine'
    )
    .action((input, options) => run(update, input, options));

  program
    .command('validate <input>')
    .description('Validate a config file')
    .action((input) => run(validate, input));

  program
    .command('plugins [input]')
    .description('List all available plugins')
    .action((input) => run(plugins, input));

  program
    .command('publish-gateway <config>')
    .alias('pg')
    .description('Publish gateway config')
    .option('--dry-run', 'Enable dry run')
    .action((input, options) => run(publish, input, options));

  program.version(pkg.version, '-v, --version');
  await program.parseAsync(process.argv);
};

main();
