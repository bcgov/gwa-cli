import React from 'react';
import { Box, Text } from 'ink';

import api from '../../services/api';
import Success from '../../components/success';
import useAsync from '../../hooks/use-async';
import type { AclContent } from './types';

interface AclResponse {
  added: number;
  missing: number;
  removed: number;
}

interface AclRequestViewProps {
  data: AclContent[];
}

const AclRequestView: React.FC<AclRequestViewProps> = ({ data }) => {
  const result = useAsync<AclResponse>(
    api,
    '/namespaces/:namespace/membership',
    {
      method: 'PUT',
      body: JSON.stringify(data),
    }
  );

  return (
    <Box flexDirection="column">
      <Success>Access Control Updated!</Success>
      <Box>
        <Box marginRight={3}>
          <Text color="green">+ Added</Text>
        </Box>
        <Text>{result.added}</Text>
      </Box>
      <Box>
        <Box marginRight={1}>
          <Text color="red">- Removed</Text>
        </Box>
        <Text>{result.removed}</Text>
      </Box>
      <Box>
        <Box marginRight={1}>
          <Text color="yellow">? Missing</Text>
        </Box>
        <Text>{result.missing}</Text>
      </Box>
    </Box>
  );
};

export default AclRequestView;
