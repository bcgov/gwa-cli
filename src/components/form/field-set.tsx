import React, { cloneElement } from 'react';
import { Box, Text, useFocus } from 'ink';

const FieldSet = ({ children, error }) => {
  const { isFocused } = useFocus();

  return (
    <Box>
      <Box width={1}>
        <Text inverse={isFocused} color="yellow">
          {isFocused ? ' ' : ' '}
        </Text>
      </Box>
      <Box width={1} marginRight={1}>
        <Text bold color="red" inverse={error}>
          {error ? '!' : ' '}
        </Text>
      </Box>
      <Box>
        {cloneElement(children, {
          focused: isFocused,
        })}
      </Box>
    </Box>
  );
};

export default FieldSet;
