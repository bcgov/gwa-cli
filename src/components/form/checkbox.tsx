import React from 'react';
import { Box, Text, useInput } from 'ink';

interface CheckboxProps {
  autoFocus: boolean;
  checked: boolean;
  error: boolean;
  focused: boolean;
  name: string;
  required: boolean;
  onChange: (key: string, value: boolean) => void;
}

const Checkbox: React.FC<CheckboxProps> = ({
  autoFocus = false,
  error,
  focused,
  checked,
  name,
  required = false,
  onChange,
}) => {
  const focusedColor = focused ? 'yellow' : 'cyan';
  const hasError = Boolean(error);
  const labelColor = hasError ? 'red' : focusedColor;

  useInput((input, key) => {
    if (focused && key.return) {
      onChange(name, !checked);
    }
  });

  return (
    <Box>
      <Box marginRight={1}>
        <Text color={labelColor}>
          {name}
          {required && '*'}:
        </Text>
      </Box>
      <Box marginRight={1}>
        <Text bold={focused}>[{checked ? 'X' : ' '}]</Text>
      </Box>
      <Box>
        <Text italic color="grey">
          Enter to toggle
        </Text>
      </Box>
    </Box>
  );
};

export default Checkbox;
