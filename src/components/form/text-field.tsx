import React, { useRef, useState, useEffect, useCallback } from 'react';
import { Box, Text, useInput } from 'ink';
import fetch from 'node-fetch';
import TextInput from 'ink-text-input';
import validUrl from 'valid-url';

interface TextFieldProps {
  encrypted: boolean;
  error?: any | undefined;
  focused?: boolean;
  name: string;
  onChange: (key: string, value: string) => void;
  onEncrypt: (key: string) => void;
  required?: boolean;
  type?: 'text' | 'url';
  value?: string | null;
}

const TextField: React.FC<TextFieldProps> = ({
  encrypted,
  error,
  focused,
  onChange,
  onEncrypt,
  name,
  required = false,
  type,
  value,
}) => {
  const [isCommand, setIsCommand] = useState<boolean>(false);
  const focusedColor = focused ? 'yellow' : 'cyan';
  const labelColor = error ? 'red' : focusedColor;
  const changeHandler = useCallback(
    (value: string) => {
      if (value.includes('~')) {
        onEncrypt(name);
      } else if (!isCommand) {
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
        <Text color={labelColor}>
          {name}
          {required && '*'}:
        </Text>
      </Box>
      <Box flexGrow={1}>
        {focused && (
          <TextInput value={value ? value : ''} onChange={changeHandler} />
        )}
        {!focused && <Text>{value === null ? 'null' : value}</Text>}
      </Box>
    </Box>
  );
};

export default TextField;
