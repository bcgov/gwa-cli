import React from 'react';
import { Box, Text } from 'ink';

interface ErrorsProps {
  errors: any;
}

const Errors: React.FC<ErrorsProps> = ({ errors }) => {
  const elements = [];

  for (const key in errors) {
    const err = errors[key];
    elements.push(
      <Box key={key}>
        <Box marginRight={3}>
          <Text bold color="red">
            {key}
          </Text>
        </Box>
        <Text color="white">{err.join(', ')}</Text>
      </Box>
    );
  }

  return (
    <Box flexDirection="column" borderColor="redBright" borderStyle="single">
      <Box>
        <Text bold color="red">
          {`Errors (${elements.length}) `}
        </Text>
      </Box>
      {elements}
    </Box>
  );
};

export default Errors;
