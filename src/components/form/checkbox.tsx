import React from 'react';
import { Box, Text, useInput } from 'ink';

interface CheckboxProps {
  checked?: boolean;
  error?: boolean;
  focused?: boolean;
  name: string;
  onChange: (key: string, value: boolean) => void;
}

const Checkbox: React.FC<CheckboxProps> = ({
  error = false,
  focused = false,
  checked,
  name,
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
        <Text color={labelColor}>{name}</Text>
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
