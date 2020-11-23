import React from 'react';
import { Box, render } from 'ink';

import AsyncAction from '../../components/async-action';
import AclRequestView from './request-view';

const renderView = (data: any, debug: boolean) => {
  return render(
    <Box flexDirection="column">
      <AsyncAction loadingText="Publishing ACL changes" verbose={debug}>
        <AclRequestView data={data} />
      </AsyncAction>
    </Box>
  );
};

export default renderView;
