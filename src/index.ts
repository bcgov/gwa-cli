import { program } from 'commander';

import run from './run';
import acl from './commands/acl';
import plugins from './commands/plugins';
import create from './commands/create';
import edit from './commands/edit';
import init from './commands/init';
import publish from './commands/publish';
import update from './commands/update';
import validate from './commands/validate';

// @ts-ignore
const pkg = require('../package.json');

const main = async () => {
  program
    .command('init')
    .option(
      '--namespace <namespace>',
      'Represents a collections of Kong Services and Routes'
    )
    .option('--service <service>', "The service's name")
    .option('--client-id <clientId>', "The service's name")
    .option('--client-secret <clientSecret>', "The service's name")
    .action((options) => run(init, null, options));

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
      'Initialize a config file in the current directory. The input file must be an OpenAPI JSON file or URL'
    )
    .action((input, options) => run(create, input, options));

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
    .option('-D, --dev', 'Dev environment')
    .option('-P, --prod', 'Production environment')
    .option('-T, --testing', 'Testing environment')
    .option('--dry-run', 'Enable dry run')
    .action((input, options) => run(publish, input, options));

  program
    .command('acl')
    .description('Update the full membership')
    .option('-D, --dev', 'Dev environment')
    .option('-P, --prod', 'Production environment')
    .option('-T, --testing', 'Testing environment')
    .option('-u, --users <users>', 'Users to add')
    .action((options) => run(acl, null, options));

  program.version(pkg.version, '-v, --version');

  await program.parseAsync(process.argv);
};

main();
