import * as React from 'react';
import { Box, Text, useApp } from 'ink';

interface FailedProps {
  error: unknown;
  verbose?: boolean;
}

const Failed: React.FC<FailedProps> = ({ error, verbose }) => {
  const { exit } = useApp();

  React.useEffect(() => {
    process.exitCode = 1;
    exit(error as Error);
  }, []);

  const e = error as { status: string; statusText: string; stack: string };

  return (
    <Box>
      <Box flexDirection="column">
        <Box>
          <Text bold color="red">
            x Error
          </Text>
          {e && (
            <Box>
              <Box marginLeft={1} flexDirection="column">
                <Text>{e.statusText}</Text>
                {e.status && (
                  <Box>
                    <Text>{`Status code ${e.status.toString()} `}</Text>
                  </Box>
                )}
              </Box>
            </Box>
          )}
        </Box>
        {verbose && e?.stack && (
          <Box marginTop={1}>
            <Text dimColor>Details</Text>
            <Box marginX={3}>
              <Text>{e.stack}</Text>
            </Box>
          </Box>
        )}
      </Box>
    </Box>
  );
};

export default Failed;
