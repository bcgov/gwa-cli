import * as React from 'react';
import { Box, Text } from 'ink';

interface ErrorViewProps {
  text: string;
  title?: string;
}

const ErrorView: React.FC<ErrorViewProps> = ({ text, title = 'Error' }) => {
  return (
    <Box borderColor="red" borderStyle="round" margin={1} paddingX={2}>
      <Box flexDirection="column">
        <Box marginBottom={1}>
          <Text underline bold color="red">
            {title}
          </Text>
        </Box>
        <Box>
          <Text dimColor>Details</Text>
          <Box marginLeft={3}>
            <Text>{text}</Text>
          </Box>
        </Box>
      </Box>
    </Box>
  );
};

export default ErrorView;
