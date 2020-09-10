import React, { cloneElement, useEffect } from 'react';
import { Box, Text, useFocus } from 'ink';

interface FieldSetProps {
  children: React.ReactElement | React.ReactElement[];
  enabled: boolean;
  editing: boolean;
  encrypted: boolean;
  error: boolean;
  focused: boolean;
  index: number;
  required: boolean;
}

const FieldSet: React.FC<FieldSetProps> = ({
  children,
  enabled,
  editing,
  error,
  encrypted,
  focused,
  index,
  required,
}) => {
  const isFocused = false;
  const focusedColor = enabled ? 'cyan' : 'yellow';
  const requiredColumn = required ? '*' : ' ';
  /* const { isFocused } = useFocus({
   *   isActive: enabled,
   * }); */

  return (
    <Box>
      <Box width={4} justifyContent="flex-end">
        <Text inverse={error || focused} color={error ? 'red' : focusedColor}>
          {encrypted ? 'E' : ' '}
        </Text>
        <Text inverse={error || focused} color={error ? 'red' : focusedColor}>
          {index < 10 && ' '} {index}
        </Text>
      </Box>
      <Box width={1} marginRight={1}>
        <Text
          bold
          inverse={error || focused}
          color={error ? 'red' : focusedColor}
        >
          {error ? '!' : requiredColumn}
        </Text>
      </Box>
      <Box flexGrow={1}>
        {cloneElement(children as React.ReactElement<any>, {
          enabled,
          editing,
          focused,
          encrypted,
        })}
      </Box>
    </Box>
  );
};

export default FieldSet;
