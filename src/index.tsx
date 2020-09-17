import * as React from 'react';
import { render } from 'ink';
import minimist, { ParsedArgs } from 'minimist';
import { Router } from 'react-router';
import { createMemoryHistory } from 'history';

import { IAppContext } from './types';
import App from './views/app';
import { loadConfig } from './services/kong';

const history = createMemoryHistory();
const args: ParsedArgs = minimist(process.argv.slice(2));
const configFile: string | null = args._[0];
const version: string = args['app-version'];

main({ dir: args.dir, file: configFile, version });

function main(args: IAppContext) {
  if (args.file) {
    // TODO: Handle if the file doesn't exist.
    loadConfig(args.file);
  }

  render(
    <Router history={history}>
      <App args={args} />
    </Router>
  );
}
