import React from 'react';
import { Box, Text } from 'ink';
import Spinner from 'ink-spinner';
import type { SpinnerName } from 'cli-spinners';

interface LoadingProps {
  children: string;
  spinner?: SpinnerName;
}

const Loading: React.FC<LoadingProps> = ({ children, spinner }) => {
  return (
    <Box>
      <Text>
        <Spinner type={spinner} /> <Text>{children}</Text>
      </Text>
    </Box>
  );
};

export default Loading;
