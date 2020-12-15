import React from 'react';
import { render } from 'ink';

import AsyncAction from '../../components/async-action';
import StatusView from './status-view';

const renderer = () => {
  render(
    <AsyncAction loadingText="Fetching status...">
      <StatusView />
    </AsyncAction>
  );
};

export default renderer;
