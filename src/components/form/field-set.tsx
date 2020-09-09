import React, { cloneElement } from 'react';
import { Box, Text, useFocus } from 'ink';

interface FieldSetProps {
  children: React.ReactElement | React.ReactElement[];
  encrypted: boolean;
  error: boolean;
  index: number;
}

const FieldSet: React.FC<FieldSetProps> = ({
  children,
  error,
  encrypted,
  index,
}) => {
  const { isFocused } = useFocus();

  return (
    <Box>
      <Box width={4} justifyContent="flex-end">
        <Text inverse={error || isFocused} color={error ? 'red' : 'yellow'}>
          {encrypted ? 'E' : ' '}
        </Text>
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
        {cloneElement(children as React.ReactElement<any>, {
          focused: isFocused,
          encrypted,
        })}
      </Box>
    </Box>
  );
};

export default FieldSet;
