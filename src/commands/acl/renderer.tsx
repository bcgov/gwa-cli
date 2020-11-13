import React from 'react';
import { Box, render } from 'ink';

import AsyncAction from '../../components/async-action';
import AclRequestView from './request-view';

const renderView = (data: any) => {
  return render(
    <Box flexDirection="column">
      <AsyncAction loadingText="Writing file...">
        <AclRequestView data={data} />
      </AsyncAction>
    </Box>
  );
};

export default renderView;
