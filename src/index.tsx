import * as React from 'react';
import { render } from 'ink';
import minimist, { ParsedArgs } from 'minimist';

import { IAppContext } from '../types';
import App from './views/app';

const args: ParsedArgs = minimist(process.argv.slice(2));
main({ dir: args.dir });

function main(args: IAppContext) {
  render(<App args={args} />);
}
