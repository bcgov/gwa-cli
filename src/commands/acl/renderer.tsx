import React from 'react';
import { Box, Instance, render } from 'ink';

import AsyncAction from '../../components/async-action';
import AclRequestView from './request-view';
import type { AclContent } from './types';

const renderView = (data: AclContent[], debug: boolean): Instance => {
  return render(
    <Box flexDirection="column">
      <AsyncAction loadingText="Publishing ACL changes" verbose={debug}>
        <AclRequestView data={data} />
      </AsyncAction>
    </Box>
  );
};

export default renderView;
