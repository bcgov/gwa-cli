import React from 'react';
import { render } from 'ink';

import CreateConfigView from './views/create-config';
import prompts from './prompts';

const renderView = () => {
  return render(<CreateConfigView prompts={prompts} />);
};

export default renderView;
