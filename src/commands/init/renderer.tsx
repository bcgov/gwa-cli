import React from 'react';
import { render } from 'ink';

import CreateEnvView from '../../views/create-env';
import prompts from './prompts';

const renderView = (env: string) => {
  return render(<CreateEnvView env={env} prompts={prompts} />);
};

export default renderView;
