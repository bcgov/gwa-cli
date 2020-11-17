import React from 'react';
import { render } from 'ink';

import AsyncAction from '../../components/async-action';
import UploadView from './upload-view';

const renderer = (options: any) => {
  render(
    <AsyncAction loadingText="Publishing gateway config...">
      <UploadView options={options} />
    </AsyncAction>
  );
};

export default renderer;
