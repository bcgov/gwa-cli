import * as React from 'react';
import TextInput from 'ink-text-input';
import { Box, Text, useInput } from 'ink';

interface NumberFieldProps {
  focused: boolean;
  name: string;
  onChange: (key: string, value: number) => void;
  required?: boolean;
  step?: number;
  min?: number;
  max?: number;
  value: number;
}

const NumberField: React.FC<NumberFieldProps> = ({
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
  const hasError = Boolean(error);
  const focusedColor = focused ? 'yellow' : '';
  const labelColor = hasError ? 'red' : focusedColor;
  const handleChange = (value: string | number) => {
    const newValue = Number(value);
    if (!Number.isNaN(newValue)) {
      onChange(name, Number(value));
    }
  };

  useInput((input, key) => {
    if (focused) {
      if (key.upArrow) {
        handleChange(value + 1);
      } else if (key.downArrow) {
        handleChange(value - 1);
      }
    }
  });

  return (
    <Box>
      <Box marginRight={1}>
        <Text bold color={labelColor}>
          {name}
          {required && '*'}:
        </Text>
      </Box>
      <Box flexGrow={1} width="50%">
        {focused && (
          <TextInput value={(value || 0).toString()} onChange={handleChange} />
        )}
        {!focused && <Text>{value}</Text>}
      </Box>
    </Box>
  );
};

export default NumberField;
