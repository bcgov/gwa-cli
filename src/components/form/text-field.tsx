import React, { useRef, useState, useEffect, useCallback } from 'react';
import TextInput from 'ink-text-input';
import { Box, Text, useInput } from 'ink';
import validUrl from 'valid-url';

interface TextFieldProps {
  error?: any | undefined;
  focused: boolean;
  name: string;
  onChange: (key: string, value: string) => void;
  required?: boolean;
  type?: 'text' | 'url';
  value?: string | null;
}

const TextField: React.FC<TextFieldProps> = ({
  error,
  focused,
  onChange,
  name,
  required = false,
  type,
  value,
}) => {
  const [isCommand, setIsCommand] = useState<boolean>(false);
  const hasError = Boolean(error);
  const focusedColor = focused ? 'yellow' : '';
  const labelColor = hasError ? 'red' : focusedColor;
  const changeHandler = useCallback(
    (value: string) => {
      if (!isCommand) {
        onChange(name, value);
      }
    },
    [focused, isCommand]
  );

  useInput((input, key) => {
    if (focused) {
      setIsCommand(key.ctrl);
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
          <TextInput value={value ? value : ''} onChange={changeHandler} />
        )}
        {!focused && <Text>{value === null ? 'null' : value}</Text>}
      </Box>
      {hasError && (
        <Box>
          <Text color="red">{error}</Text>
        </Box>
      )}
    </Box>
  );
};

export default TextField;
