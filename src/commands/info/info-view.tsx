import * as React from 'react';
import { Box, Text, Transform } from 'ink';

import api from '../../services/api';
import useAsync from '../../hooks/use-async';
import type { InfoData } from './types';
import config from '../../config';

interface InfoViewProps {}

const InfoView: React.FC<InfoViewProps> = () => {
  const data = useAsync<InfoData>(api, '/whoami');
  const { status } = useAsync<InfoData>(api, '/status');
  const { env, clientSecret } = config();
  const isUp = status === 'ok';
  const statusTextColor = isUp ? 'green' : 'red';
  const statusIcon = isUp ? '✓' : 'x';
  const labelColumnWidth = 15;

  return (
    <Box flexDirection="column" marginY={1}>
      <Box marginBottom={1}>
        <Text>
          Namespace: <Text bold>{data.namespace}</Text>
        </Text>
      </Box>
      <Box>
        <Box width={labelColumnWidth}>
          <Text color="cyan">Environment</Text>
        </Box>
        <Box>
          <Text>{env}</Text>
        </Box>
      </Box>
      <Box>
        <Box width={labelColumnWidth}>
          <Text color="cyan">Client ID</Text>
        </Box>
        <Box>
          <Text>{data.authorizedParty}</Text>
        </Box>
      </Box>
      <Box>
        <Box width={labelColumnWidth}>
          <Text color="cyan">Client Secret</Text>
        </Box>
        <Box>
          <Transform transform={(output) => '*'.repeat(output.length)}>
            <Text>{clientSecret}</Text>
          </Transform>
        </Box>
      </Box>
      <Box>
        <Box width={labelColumnWidth}>
          <Text color="cyan">API Status</Text>
        </Box>
        <Box>
          <Text color={statusTextColor}>{`${statusIcon} ${status}`}</Text>
        </Box>
      </Box>
    </Box>
  );
};

export default InfoView;
