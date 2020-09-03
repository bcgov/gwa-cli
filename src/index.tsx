import * as React from 'react';
import { render } from 'ink';
import minimist, { ParsedArgs } from 'minimist';
import { Router, Route, useHistory } from 'react-router';
import { createMemoryHistory } from 'history';

import { IAppContext } from './types';
import App from './views/app';

const history = createMemoryHistory();
const args: ParsedArgs = minimist(process.argv.slice(2));
main({ dir: args.dir });

function main(args: IAppContext) {
  render(
    <Router history={history}>
      <App args={args} />
    </Router>
  );
}
