import React from 'react';
import { Box, Text } from 'ink';
import Spinner from 'ink-spinner';

interface HttpRequestProps {
  loading?: boolean;
  loadingText?: string;
  successText?: string;
}

const HttpRequest: React.FC<HttpRequestProps> = ({
  loading,
  loadingText,
  successText,
}) => {
  if (loading) {
    return (
      <Box>
        <Text>
          <Spinner /> <Text>{loadingText}</Text>
        </Text>
      </Box>
    );
  }

  return (
    <Box>
      <Box marginRight={1}>
        <Text bold color="green">
          ✓
        </Text>
      </Box>
      <Text bold>
        <Text>{successText}</Text>
      </Text>
    </Box>
  );
};

export default HttpRequest;
