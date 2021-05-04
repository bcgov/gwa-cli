import React from 'react';
import { render } from 'ink';

import AsyncAction from '../../components/async-action';
import UploadView from './upload-view';

interface PublishRenderOptions {
  configFile?: string;
  dryRun: string;
  verbose?: boolean;
}
const renderer = (options: PublishRenderOptions) => {
  render(
    <AsyncAction
      loadingText="Publishing gateway config..."
      verbose={options.verbose}
    >
      <UploadView options={options} />
    </AsyncAction>
  );
};

export default renderer;
