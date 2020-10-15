import React from 'react';
import { Box, Text } from 'ink';

interface SuccessProps {
  children: string | React.ReactNode;
}

const Success: React.FC<SuccessProps> = ({ children }) => {
  return (
    <Box>
      <Box marginRight={1}>
        <Text bold color="green">
          ✓
        </Text>
      </Box>
      <Text bold>
        <Text>{children}</Text>
      </Text>
    </Box>
  );
};

export default Success;
