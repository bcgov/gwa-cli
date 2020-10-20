import * as React from 'react';
import { render } from 'ink';
import { Router } from 'react-router';
import { createMemoryHistory } from 'history';

import { IAppContext } from './types';
import App from './views/app';

const history = createMemoryHistory();

export default function (args: IAppContext) {
  render(
    <Router history={history}>
      <App args={args} />
    </Router>
  );
}
