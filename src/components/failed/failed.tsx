import * as React from 'react';
import { Box, Text } from 'ink';
import type { FallbackProps } from 'react-error-boundary';

interface FailedProps extends FallbackProps {}

const Failed: React.FC<FailedProps> = ({ error }) => {
  return (
    <Box>
      <Box flexDirection="column">
        <Box marginBottom={1}>
          <Text bold color="red">
            x Action Failed
          </Text>
        </Box>
        {error?.message && (
          <Box>
            <Text dimColor>Details</Text>
            <Box marginX={3}>
              <Text>{error.message}</Text>
            </Box>
          </Box>
        )}
      </Box>
    </Box>
  );
};

export default Failed;
