import React from 'react';
import { Box, Text, useFocus, useInput } from 'ink';

interface CheckboxProps {
  autoFocus: boolean;
  checked: boolean;
  label: string;
  name: string;
  onChange: (value: boolean) => void;
}

const Checkbox: React.FC<CheckboxProps> = ({
  autoFocus = false,
  checked,
  label,
  name,
  onChange,
}) => {
  const { isFocused } = useFocus({ autoFocus });

  useInput((input, key) => {
    if (isFocused && key.return) {
      onChange(!checked);
    }
  });

  return (
    <Box>
      <Box marginRight={1}>
        <Text bold={isFocused}>[{checked ? 'X' : ' '}]</Text>
      </Box>
      <Box marginRight={3}>
        <Text>{label}</Text>
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
