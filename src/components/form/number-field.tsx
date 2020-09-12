import React, { useCallback } from 'react';
import { Box, Text, useInput } from 'ink';

interface NumberFieldProps {
  enabled?: boolean;
  focused?: boolean;
  name: string;
  onChange: (key: string, value: number) => void;
  required?: boolean;
  step?: number;
  min?: number;
  max?: number;
  value: number;
}

const NumberField: React.FC<NumberFieldProps> = ({
  enabled = false,
  focused,
  onChange,
  max,
  min,
  name,
  required = false,
  step,
  value = 0,
}) => {
  const [error, setError] = React.useState<string>('');
  const stringValue = value ? value.toString() : '';
  const hasError = Boolean(error);
  const focusedColor = focused ? 'yellow' : 'cyan';
  const labelColor = hasError ? 'red' : focusedColor;
  const handleChange = useCallback((val: string | number) => {
    const newValue = Number(val);

    if (!Number.isNaN(newValue)) {
      onChange(name, newValue);
    }
  }, []);

  useInput((input, key) => {
    if (focused && enabled) {
      if (key.upArrow) {
        handleChange(value + 1);
      } else if (key.downArrow) {
        handleChange(value - 1);
      } else if (key.delete) {
        handleChange(stringValue.slice(0, stringValue.length - 1));
      } else if (!Number.isNaN(input)) {
        handleChange(stringValue + input);
      }
    }
  });

  return (
    <Box>
      <Box marginRight={1}>
        <Text color={labelColor}>{name}</Text>
      </Box>
      <Box flexGrow={1} width="50%">
        <Text>{(!value ? 0 : value).toString()}</Text>
      </Box>
    </Box>
  );
};

export default NumberField;
