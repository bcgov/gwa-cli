import React from 'react';
import { render } from 'ink';

import AsyncAction from '../../components/async-action';
import InfoView from './info-view';

const renderer = () => {
  render(
    <AsyncAction loadingText="Fetching your details...">
      <InfoView />
    </AsyncAction>
  );
};

export default renderer;
