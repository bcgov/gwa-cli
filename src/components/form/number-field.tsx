import * as React from 'react';
import TextInput from 'ink-text-input';
import { Box, Text, useFocus } from 'ink';

interface NumberFieldProps {
  name: string;
  onChange: (key: string, value: number) => void;
  required?: boolean;
  step?: number;
  min?: number;
  max?: number;
  value: number;
}

function NumberField({
  onChange,
  max,
  min,
  name,
  required = false,
  step,
  value = 0,
}: NumberFieldProps) {
  const hasEntered = React.useRef<boolean>(false);
  const [error, setError] = React.useState<string>('');
  const { isFocused } = useFocus();
  const hasError = Boolean(error) && hasEntered.current;
  const focusedColor = isFocused ? 'green' : '';
  const labelColor = hasError ? 'red' : focusedColor;
  const handleChange = (value: string) => {
    hasEntered.current = true;
    onChange(name, Number(value));
  };

  /* React.useEffect(() => {
   *   if (!isFocused && required && !value) {
   *     setError('This field is required');
   *   } else {
   *     setError('');
   *   }
   * }, [isFocused, required, value]); */

  return (
    <Box>
      <Box marginRight={1}>
        <Text bold color={labelColor}>
          {name}
          {required && '*'}:
        </Text>
      </Box>
      <Box flexGrow={1} width="50%">
        {isFocused && (
          <TextInput value={value.toString()} onChange={handleChange} />
        )}
        {!isFocused && <Text>{value}</Text>}
      </Box>
      {hasError && (
        <Box>
          <Text color="red">{error}</Text>
        </Box>
      )}
    </Box>
  );
}

export default NumberField;
