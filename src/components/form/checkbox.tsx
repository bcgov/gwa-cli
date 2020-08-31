import React from 'react';
import { Box, Text, useFocus, useInput } from 'ink';

interface CheckboxProps {
  checked: boolean;
  label: string;
  name: string;
  onChange: (value: boolean) => void;
}

const Checkbox: React.FC<CheckboxProps> = ({
  checked,
  label,
  name,
  onChange,
}) => {
  const { isFocused } = useFocus();

  useInput((input, key) => {
    if (isFocused && key.return) {
      onChange(!checked);
    }
  });

  return (
    <Box>
      <Text bold={isFocused}>
        <Text>[{checked ? 'X' : ' '}]</Text> {label}
      </Text>
    </Box>
  );
};

export default Checkbox;
