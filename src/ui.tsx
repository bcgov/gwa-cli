import * as React from 'react';
import { render } from 'ink';
import { Router, Redirect } from 'react-router';
import { createMemoryHistory } from 'history';

import { IAppContext } from './types';
import App from './views/app';

const history = createMemoryHistory();

export default function (args: IAppContext, redirect: string = '/') {
  render(
    <Router history={history}>
      <Redirect from="/" to={redirect} />
      <App args={args} />
    </Router>
  );
}
