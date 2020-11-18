import React from 'react';
import { render } from 'ink';

import CreateConfigView from './views/create-config';
import prompts from './prompts';

const renderView = (submitHandler: any) => {
  return render(
    <CreateConfigView prompts={prompts} submitHandler={submitHandler} />
  );
};

export default renderView;
