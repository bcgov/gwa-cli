import React from 'react';
import { Box, Text } from 'ink';
import Spinner from 'ink-spinner';
import type { SpinnerName } from 'cli-spinners';

interface LoadingProps {
  children?: React.ReactNode;
  spinner?: SpinnerName;
  text?: string;
}

const Loading: React.FC<LoadingProps> = ({
  children,
  spinner,
  text = 'Processing...',
}) => {
  return (
    <Box flexDirection="column">
      <Text>
        <Spinner type={spinner} /> <Text>{text}</Text>
      </Text>
      {children && <Box>{children}</Box>}
    </Box>
  );
};

export default Loading;
