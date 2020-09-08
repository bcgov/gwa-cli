import React, { cloneElement } from 'react';
import { Box, Text, useFocus } from 'ink';

interface FieldSetProps {
  children: React.ReactElement | React.ReactElement[];
  error: boolean;
  index: number;
}

const FieldSet: React.FC<FieldSetProps> = ({ children, error, index }) => {
  const { isFocused } = useFocus();

  return (
    <Box>
      <Box width={3} justifyContent="flex-end">
        <Text inverse={error || isFocused} color={error ? 'red' : 'yellow'}>
          {index < 10 && ' '} {index}
        </Text>
      </Box>
      <Box width={1} marginRight={1}>
        <Text
          bold
          inverse={error || isFocused}
          color={error ? 'red' : 'yellow'}
        >
          {error ? '!' : ' '}
        </Text>
      </Box>
      <Box flexGrow={1}>
        {cloneElement(children, {
          focused: isFocused,
        })}
      </Box>
    </Box>
  );
};

export default FieldSet;
