import * as React from 'react';
import { render } from 'ink';
import { Router, Redirect } from 'react-router';
import { createMemoryHistory } from 'history';

import App from './views/app';

const history = createMemoryHistory();

export default function (redirect: string = '/') {
  render(
    <Router history={history}>
      {redirect && <Redirect from="/" to={redirect} />}
      <App />
    </Router>
  );
}
