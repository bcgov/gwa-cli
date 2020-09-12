import React, { cloneElement, useEffect } from 'react';
import { Box, Text, useInput } from 'ink';

interface FieldSetProps {
  children: React.ReactElement | React.ReactElement[];
  enabled: boolean;
  editing: boolean;
  encrypted: boolean;
  error: boolean;
  focused: boolean;
  onEncrypt: (name: string, encrypted: boolean) => void;
  name: string;
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
  name,
  onEncrypt,
  required,
}) => {
  const focusedColor = enabled ? 'cyan' : 'yellow';
  const indicator = focused ? '>' : ' ';
  const requiredColumn = required ? '*' : ' ';

  useInput((input) => {
    if (focused && input === 'E') {
      onEncrypt(name, !encrypted);
    }
  });

  return (
    <Box>
      <Box width={5}>
        <Text inverse={error || focused} color={error ? 'red' : focusedColor}>
          {indicator}
        </Text>
        <Text inverse={error || focused} color={error ? 'red' : focusedColor}>
          {encrypted ? 'E ' : '  '}
          {index < 10 && ' '}
          {index}
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
