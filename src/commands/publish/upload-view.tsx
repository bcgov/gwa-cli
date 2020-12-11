import React from 'react';
import { Box, Text } from 'ink';

import publish from '../../services/publish';
import Success from '../../components/success';
import useAsync from '../../hooks/use-async';

interface UploadViewProps {
  options: {
    input?: string;
    dryRun: string;
  };
}

const UploadView: React.FC<UploadViewProps> = ({ options }) => {
  const { message, results } = useAsync(
    publish,
    '/namespaces/:namespace/gateway',
    options
  );

  return (
    <Box flexDirection="column">
      <Success>
        <Text color="green">Success</Text>
        {` Configuration ${options.input} Published`}
      </Success>
      <Box flexDirection="column" marginTop={1}>
        <Box>
          <Text>{message}</Text>
        </Box>
        <Box marginTop={1}>
          <Text>{results}</Text>
        </Box>
      </Box>
    </Box>
  );
};

export default UploadView;
