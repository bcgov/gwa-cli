import * as React from 'react';
import { Box, Text, useApp } from 'ink';
import type { FallbackProps } from 'react-error-boundary';

interface FailedProps {
  error: Error | undefined;
  verbose: boolean;
}

const Failed: React.FC<FailedProps> = ({ error, verbose }) => {
  const { exit } = useApp();

  React.useEffect(() => {
    process.exitCode = 1;
    exit(error);
  }, []);

  return (
    <Box>
      <Box flexDirection="column">
        <Box>
          <Text bold color="red">
            x Error
          </Text>
          {error && (
            <Box marginLeft={1}>
              <Text>{error.message}</Text>
            </Box>
          )}
        </Box>
        {verbose && error?.stack && (
          <Box marginTop={1}>
            <Text dimColor>Details</Text>
            <Box marginX={3}>
              <Text>{error.stack}</Text>
            </Box>
          </Box>
        )}
      </Box>
    </Box>
  );
};

export default Failed;
